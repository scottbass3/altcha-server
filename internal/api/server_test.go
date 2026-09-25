package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/altcha-org/altcha-lib-go"
	"github.com/scottbass3/altcha-server/internal/config"
	"github.com/scottbass3/altcha-server/internal/store"
)

const testKey = "test-hmac-key-0123456789"

func newTestHandler(t *testing.T, mutate func(*config.Config)) http.Handler {
	t.Helper()
	cfg := config.Config{
		HmacKey:     testKey,
		MaxNumber:   1000,
		Algorithm:   "SHA-256",
		SaltLength:  12,
		Expire:      "10m",
		CheckExpire: true,
		CorsOrigins: "*",
	}
	if mutate != nil {
		mutate(&cfg)
	}
	s, err := NewServer(cfg)
	if err != nil {
		t.Fatal(err)
	}
	s.nonces = store.NewMemoryStore(t.Context())
	return s.routes()
}

func do(t *testing.T, h http.Handler, method, path, body string, headers ...string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	for i := 0; i+1 < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func mustJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func solvedPayload(t *testing.T, h http.Handler) altcha.Payload {
	t.Helper()
	rec := do(t, h, "GET", "/request", "")
	var ch altcha.Challenge
	if err := json.Unmarshal(rec.Body.Bytes(), &ch); err != nil {
		t.Fatalf("decode challenge: %v (%s)", err, rec.Body)
	}
	sol, err := altcha.SolveChallenge(ch.Challenge, ch.Salt, altcha.Algorithm(ch.Algorithm), int(ch.MaxNumber), 0, nil)
	if err != nil || sol == nil {
		t.Fatalf("solve failed: %v", err)
	}
	return altcha.Payload{Algorithm: ch.Algorithm, Challenge: ch.Challenge, Number: int64(sol.Number), Salt: ch.Salt, Signature: ch.Signature}
}

func signedPayload(t *testing.T, number int64, expires time.Time) altcha.Payload {
	t.Helper()
	ch, err := altcha.CreateChallenge(altcha.ChallengeOptions{HMACKey: testKey, Number: &number, Expires: &expires})
	if err != nil {
		t.Fatal(err)
	}
	return altcha.Payload{Algorithm: ch.Algorithm, Challenge: ch.Challenge, Number: number, Salt: ch.Salt, Signature: ch.Signature}
}

func TestHealth(t *testing.T) {
	if rec := do(t, newTestHandler(t, nil), "GET", "/health", ""); rec.Code != http.StatusOK {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestBaseURL(t *testing.T) {
	h := newTestHandler(t, func(c *config.Config) { c.BaseUrl = "/altcha" })
	if rec := do(t, h, "GET", "/altcha/request", ""); rec.Code != http.StatusOK {
		t.Errorf("prefixed route status = %d", rec.Code)
	}
	if rec := do(t, h, "GET", "/request", ""); rec.Code != http.StatusNotFound {
		t.Errorf("unprefixed route status = %d", rec.Code)
	}
}

func TestVerify(t *testing.T) {
	h := newTestHandler(t, nil)
	p := solvedPayload(t, h)

	rec := do(t, h, "POST", "/verify", mustJSON(t, p))
	if rec.Code != http.StatusOK {
		t.Fatalf("valid solution: status = %d, body = %s", rec.Code, rec.Body)
	}
	var resp struct {
		Success bool           `json:"success"`
		Data    altcha.Payload `json:"data"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if !resp.Success || resp.Data != p {
		t.Errorf("unexpected body %s", rec.Body)
	}

	if rec := do(t, h, "POST", "/verify", mustJSON(t, p)); rec.Code != http.StatusConflict {
		t.Errorf("replay: status = %d, want 409", rec.Code)
	}

	bad := solvedPayload(t, h)
	bad.Number++
	if rec := do(t, h, "POST", "/verify", mustJSON(t, bad)); rec.Code != http.StatusBadRequest {
		t.Errorf("wrong number: status = %d, want 400", rec.Code)
	}

	if rec := do(t, h, "POST", "/verify", "{not json"); rec.Code != http.StatusBadRequest {
		t.Errorf("malformed JSON: status = %d, want 400", rec.Code)
	}
}

func TestVerifyRejectsExpiredAndSpliced(t *testing.T) {
	h := newTestHandler(t, nil)
	p := signedPayload(t, 12345, time.Now().Add(-time.Hour))
	if rec := do(t, h, "POST", "/verify", mustJSON(t, p)); rec.Code != http.StatusBadRequest {
		t.Errorf("expired: status = %d, want 400", rec.Code)
	}

	spliced := p
	spliced.Salt = strings.TrimSuffix(p.Salt, "&") + "1"
	spliced.Number = 2345
	if rec := do(t, h, "POST", "/verify", mustJSON(t, spliced)); rec.Code != http.StatusBadRequest {
		t.Errorf("spliced: status = %d, want 400", rec.Code)
	}
}

func TestVerifyCheckExpireDisabled(t *testing.T) {
	h := newTestHandler(t, func(c *config.Config) { c.CheckExpire = false })
	p := signedPayload(t, 42, time.Now().Add(-time.Hour))
	if rec := do(t, h, "POST", "/verify", mustJSON(t, p)); rec.Code != http.StatusOK {
		t.Errorf("first submit: status = %d, body = %s", rec.Code, rec.Body)
	}
	if rec := do(t, h, "POST", "/verify", mustJSON(t, p)); rec.Code != http.StatusConflict {
		t.Errorf("replay: status = %d, want 409", rec.Code)
	}
}

func TestNonceExpiryIsCapped(t *testing.T) {
	s, err := NewServer(config.Config{HmacKey: testKey, Algorithm: "SHA-256", Expire: "10m", CheckExpire: true})
	if err != nil {
		t.Fatal(err)
	}
	limit := time.Now().Add(10 * time.Minute)

	forged := altcha.Payload{Salt: fmt.Sprintf("abc?expires=%d&", time.Now().Add(1000*time.Hour).Unix())}
	if got := s.nonceExpiry(forged); got.After(limit.Add(time.Second)) {
		t.Errorf("forged expiry not capped: %v", got)
	}

	soon := time.Now().Add(time.Minute)
	legit := altcha.Payload{Salt: fmt.Sprintf("abc?expires=%d&", soon.Unix())}
	if got := s.nonceExpiry(legit); got.Before(soon) {
		t.Errorf("nonce expiry %v before challenge expiry %v", got, soon)
	}
}

func TestBodyTooLarge(t *testing.T) {
	body := `{"challenge":"` + strings.Repeat("a", 1<<20) + `"}`
	if rec := do(t, newTestHandler(t, nil), "POST", "/verify", body); rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want 413", rec.Code)
	}
}

func TestCORS(t *testing.T) {
	tests := []struct {
		origins, origin, want string
	}{
		{"*", "https://a.test", "*"},
		{"https://a.test, https://b.test", "https://b.test", "https://b.test"},
		{"https://a.test, https://b.test", "https://c.test", ""},
	}
	for _, tt := range tests {
		h := newTestHandler(t, func(c *config.Config) { c.CorsOrigins = tt.origins })
		rec := do(t, h, "OPTIONS", "/verify", "", "Origin", tt.origin)
		if rec.Code != http.StatusOK {
			t.Errorf("%q: preflight status = %d", tt.origins, rec.Code)
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != tt.want {
			t.Errorf("%q with Origin %q: ACAO = %q, want %q", tt.origins, tt.origin, got, tt.want)
		}
	}
}

func TestVerifyFields(t *testing.T) {
	h := newTestHandler(t, nil)
	sum := sha256.Sum256([]byte("a@b.c\nhi"))
	req := map[string]any{
		"formData": map[string][]string{"email": {"a@b.c"}, "msg": {"hi"}},
		"fields":   []string{"email", "msg"},
	}

	req["fieldsHash"] = hex.EncodeToString(sum[:])
	if rec := do(t, h, "POST", "/verify-fields", mustJSON(t, req)); rec.Body.String() != `{"success":true}` {
		t.Errorf("good hash: %d %s", rec.Code, rec.Body)
	}
	req["fieldsHash"] = "00"
	if rec := do(t, h, "POST", "/verify-fields", mustJSON(t, req)); rec.Body.String() != `{"success":false}` {
		t.Errorf("bad hash: %d %s", rec.Code, rec.Body)
	}
}

func TestVerifyServerSignature(t *testing.T) {
	h := newTestHandler(t, nil)
	data := fmt.Sprintf("verified=true&expire=%d&score=0.5", time.Now().Add(time.Hour).Unix())
	sum := sha256.Sum256([]byte(data))
	mac := hmac.New(sha256.New, []byte(testKey))
	mac.Write(sum[:])

	for sig, want := range map[string]bool{hex.EncodeToString(mac.Sum(nil)): true, "bad": false} {
		body := mustJSON(t, map[string]any{"algorithm": "SHA-256", "verificationData": data, "signature": sig, "verified": true})
		rec := do(t, h, "POST", "/verify-server-signature", body)
		var resp struct {
			Success bool `json:"success"`
		}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		if rec.Code != http.StatusOK || resp.Success != want {
			t.Errorf("signature %q: %d %s, want success=%v", sig, rec.Code, rec.Body, want)
		}
	}
}
