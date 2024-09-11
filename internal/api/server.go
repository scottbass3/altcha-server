package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"forge.cadoles.com/cadoles/altcha-server/internal/client"
	"forge.cadoles.com/cadoles/altcha-server/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"gitlab.com/wpetit/goweb/logger"
)

type Server struct {
	baseUrl	string
	port	string
	client	client.Client
}

func (s *Server) Run(ctx context.Context) {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})
	r.Get(s.baseUrl+"/request", s.requestHandler)
	r.Get(s.baseUrl+"/verify", s.submitHandler)
	r.Get(s.baseUrl+"/verify-spam-filter", s.submitSpamFilterHandler)
	
	logger.Info(ctx, "altcha server listening on port "+s.port)
	if err := http.ListenAndServe(":"+s.port, r); err != nil {
		logger.Error(ctx, err.Error())
	}
}

func (s *Server) requestHandler(w http.ResponseWriter, r *http.Request) {
	challenge, err := s.client.Generate()

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create challenge : %s", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, challenge)
}

func (s *Server) submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	formData := r.FormValue("altcha")
	if formData == "" {
		 http.Error(w, "Atlcha payload missing", http.StatusBadRequest)
		 return
	}

	decodedPayload, err := base64.StdEncoding.DecodeString(formData)
	if err != nil {
		http.Error(w, "Failed to decode Altcha payload", http.StatusBadRequest)
		return
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(decodedPayload, &payload); err != nil {
		http.Error(w, "Failed to parse Altcha payload", http.StatusBadRequest)
		return
	}

	verified, err := s.client.VerifySolution(payload)
	
	if err != nil || !verified {
		http.Error(w, "Invalid Altcha payload", http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]interface{}{
		"success":	true,
		"data":		formData,
	})
}

func (s *Server) submitSpamFilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	formData, err := formToMap(r)
	if err != nil {
		http.Error(w, "Cannot read form data", http.StatusBadRequest)
	}

	payload := r.FormValue("altcha")
	if payload == "" {
		http.Error(w, "Atlcha payload missing", http.StatusBadRequest)
	}

	verified, verificationData, err := s.client.VerifyServerSignature(payload)
	if err != nil || !verified {
		http.Error(w, "Invalid Altcha payload", http.StatusBadRequest)
		return
	}

	if verificationData.Verified && verificationData.Expire > time.Now().Unix() {
		if verificationData.Classification == "BAD" {
			http.Error(w, "Classified as spam", http.StatusBadRequest)
			return
		}

		if verificationData.FieldsHash != "" {
			verified, err := s.client.VerifyFieldsHash(formData, verificationData.Fields, verificationData.FieldsHash)
			if err != nil || !verified {
				http.Error(w, "Invalid fields hash", http.StatusBadRequest)
				return
			}
		}

		writeJSON(w, map[string]interface{}{
			"success":			true,
			"data":				formData,
			"verificationData":	verificationData,
		})
		return
	}

	http.Error(w, "Invalid Altcha payload", http.StatusBadRequest)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func formToMap(r *http.Request) (map[string][]string, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	return r.Form, nil
}

func NewServer(cfg config.Config) *Server {
	client := *client.NewClient(cfg.HmacKey, cfg.MaxNumber, cfg.Algorithm, cfg.Salt, cfg.Expire, cfg.CheckExpire)

	return &Server {
		baseUrl: cfg.BaseUrl,
		port:	cfg.Port,
		client:	client,
	}
}