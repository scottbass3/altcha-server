package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/altcha-org/altcha-lib-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/scottbass3/altcha-server/internal/client"
	"github.com/scottbass3/altcha-server/internal/config"
	"github.com/scottbass3/altcha-server/internal/store"
	"gitlab.com/wpetit/goweb/logger"
)

type Server struct {
	baseUrl string
	port    string
	client  client.Client
	config  config.Config
	nonces  *store.NonceStore
}

func (s *Server) Run(ctx context.Context) {
	s.nonces = store.NewNonceStore(ctx)
	if s.config.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(s.corsMiddleware)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get(s.baseUrl+"/request", s.requestHandler)

	r.Post(s.baseUrl+"/verify", s.submitHandler)
	r.Post(s.baseUrl+"/verify-fields", s.verifyFieldsHandler)
	r.Post(s.baseUrl+"/verify-server-signature", s.verifyServerSignatureHandler)

	logger.Info(ctx, "altcha server listening on port "+s.port)
	if err := http.ListenAndServe(":"+s.port, r); err != nil {
		logger.Error(ctx, err.Error())
	}
}

func (s *Server) requestHandler(w http.ResponseWriter, r *http.Request) {
	challenge, err := s.client.Generate()
	if err != nil {
		slog.Debug("Failed to create challenge", "error", err)
		http.Error(w, fmt.Sprintf("Failed to create challenge: %s", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, challenge)
}

func (s *Server) submitHandler(w http.ResponseWriter, r *http.Request) {
	var payload altcha.Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		slog.Debug("Failed to parse Altcha payload", "error", err)
		http.Error(w, "Failed to parse Altcha payload", http.StatusBadRequest)
		return
	}

	verified, err := s.client.VerifySolution(payload)
	if err != nil {
		slog.Debug("Invalid Altcha payload", "error", err)
		http.Error(w, "Invalid Altcha payload", http.StatusBadRequest)
		return
	}

	if !verified && !s.config.DisableValidation {
		slog.Debug("Invalid solution")
		http.Error(w, "Invalid solution", http.StatusBadRequest)
		return
	}

	if verified && !s.nonces.Consume(payload.Challenge, challengeExpiry(payload)) {
		http.Error(w, "Challenge already used", http.StatusConflict)
		return
	}

	writeJSON(w, map[string]bool{"success": true})
}

func challengeExpiry(payload altcha.Payload) time.Time {
	params := altcha.ExtractParams(payload)
	if exp := params.Get("expires"); exp != "" {
		if unix, err := strconv.ParseInt(exp, 10, 64); err == nil {
			return time.Unix(unix, 0)
		}
	}
	return time.Now().Add(24 * time.Hour)
}

type verifyFieldsRequest struct {
	FormData   map[string][]string `json:"formData"`
	Fields     []string            `json:"fields"`
	FieldsHash string              `json:"fieldsHash"`
}

func (s *Server) verifyFieldsHandler(w http.ResponseWriter, r *http.Request) {
	var req verifyFieldsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to parse request", http.StatusBadRequest)
		return
	}

	verified, err := s.client.VerifyFieldsHash(req.FormData, req.Fields, req.FieldsHash)
	if err != nil {
		slog.Debug("VerifyFieldsHash error", "error", err)
		http.Error(w, "Failed to verify fields hash", http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]bool{"success": verified})
}

func (s *Server) verifyServerSignatureHandler(w http.ResponseWriter, r *http.Request) {
	var payload altcha.ServerSignaturePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Failed to parse request", http.StatusBadRequest)
		return
	}

	verified, data, err := s.client.VerifyServerSignature(payload)
	if err != nil {
		slog.Debug("VerifyServerSignature error", "error", err)
		http.Error(w, "Failed to verify server signature", http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]any{"success": verified, "data": data})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", s.config.CorsOrigins)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func NewServer(cfg config.Config) (*Server, error) {
	expirationDuration, err := time.ParseDuration(cfg.Expire)
	if err != nil {
		return nil, fmt.Errorf("invalid ALTCHA_EXPIRE value %q: %w", cfg.Expire, err)
	}

	c, err := client.New(cfg.HmacKey, cfg.MaxNumber, cfg.Algorithm, cfg.Salt, cfg.SaltLength, expirationDuration, cfg.CheckExpire)
	if err != nil {
		return nil, err
	}

	return &Server{
		baseUrl: cfg.BaseUrl,
		port:    cfg.Port,
		client:  *c,
		config:  cfg,
	}, nil
}
