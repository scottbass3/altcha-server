package client

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/altcha-org/altcha-lib-go"
)

const testKey = "test-hmac-key-0123456789"

func newTestClient(t *testing.T) *Client {
	t.Helper()
	c, err := New(testKey, 1000, "SHA-256", "", 12, 10*time.Minute, true)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestNewValidation(t *testing.T) {
	if _, err := New("short", 1000, "SHA-256", "", 12, time.Minute, true); err == nil {
		t.Error("expected error for short HMAC key")
	}
	if _, err := New(testKey, 1000, "MD5", "", 12, time.Minute, true); err == nil {
		t.Error("expected error for unsupported algorithm")
	}
}

func TestGenerateSolveVerify(t *testing.T) {
	c := newTestClient(t)
	ch, err := c.Generate()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(ch.Salt, "?expires=") {
		t.Errorf("salt %q has no expires param", ch.Salt)
	}

	sol, err := altcha.SolveChallenge(ch.Challenge, ch.Salt, altcha.SHA256, int(ch.MaxNumber), 0, nil)
	if err != nil || sol == nil {
		t.Fatalf("solve failed: %v", err)
	}
	payload := altcha.Payload{Algorithm: ch.Algorithm, Challenge: ch.Challenge, Number: int64(sol.Number), Salt: ch.Salt, Signature: ch.Signature}

	if ok, err := c.VerifySolution(payload); err != nil || !ok {
		t.Errorf("valid solution rejected: %v, %v", ok, err)
	}
	payload.Number++
	if ok, _ := c.VerifySolution(payload); ok {
		t.Error("wrong number accepted")
	}
}

// expiredPayload returns a correctly signed payload whose challenge expired an hour ago.
func expiredPayload(t *testing.T) altcha.Payload {
	t.Helper()
	number := int64(12345)
	expires := time.Now().Add(-time.Hour)
	ch, err := altcha.CreateChallenge(altcha.ChallengeOptions{HMACKey: testKey, Number: &number, Expires: &expires})
	if err != nil {
		t.Fatal(err)
	}
	return altcha.Payload{Algorithm: ch.Algorithm, Challenge: ch.Challenge, Number: number, Salt: ch.Salt, Signature: ch.Signature}
}

func TestVerifyRejectsExpired(t *testing.T) {
	if ok, _ := newTestClient(t).VerifySolution(expiredPayload(t)); ok {
		t.Error("expired challenge accepted")
	}
}

// GO-2025-4239: moving digits from the number into the expires param must not extend expiry.
func TestVerifyRejectsSplicedExpiry(t *testing.T) {
	p := expiredPayload(t)
	n := fmt.Sprint(p.Number)
	spliced := p
	spliced.Salt = strings.TrimSuffix(p.Salt, "&") + n[:1]
	fmt.Sscan(n[1:], &spliced.Number)

	if ok, _ := newTestClient(t).VerifySolution(spliced); ok {
		t.Errorf("spliced payload accepted: salt=%q number=%d", spliced.Salt, spliced.Number)
	}
}
