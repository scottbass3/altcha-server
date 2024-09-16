package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"forge.cadoles.com/cadoles/altcha-server/internal/client"
	"forge.cadoles.com/cadoles/altcha-server/internal/config"
	"github.com/altcha-org/altcha-lib-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"gitlab.com/wpetit/goweb/logger"
)

type Server struct {
	baseUrl	string
	port	string
	client	client.Client
	config	config.Config
}

func (s *Server) Run(ctx context.Context) {
	if s.config.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})
	r.Get(s.baseUrl+"/request", s.requestHandler)
	r.Post(s.baseUrl+"/verify", s.submitHandler)

	logger.Info(ctx, "altcha server listening on port "+s.port)
	if err := http.ListenAndServe(":"+s.port, r); err != nil {
		logger.Error(ctx, err.Error())
	}
}

func (s *Server) requestHandler(w http.ResponseWriter, r *http.Request) {
	challenge, err := s.client.Generate()

	if err != nil {
		slog.Debug("Failed to create challenge,", "error", err)
		http.Error(w, fmt.Sprintf("Failed to create challenge : %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(challenge)
	if err != nil {
		slog.Debug("Failed to encode JSON", "error", err)
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

func (s *Server) submitHandler(w http.ResponseWriter, r *http.Request) {
	var payload altcha.Payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		slog.Debug("Failed to parse Altcha payload,", "error", err)
		http.Error(w, "Failed to parse Altcha payload", http.StatusBadRequest)
		return
	}

	verified, err := s.client.VerifySolution(payload)
	
	if err != nil {
		slog.Debug("Invalid Altcha payload", "error", err)
		http.Error(w, "Invalid Altcha payload,", http.StatusBadRequest)
		return
	}

	if !verified {
		slog.Debug("Invalid solution")
		http.Error(w, "Invalid solution", http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"success":	true,
		"data":		payload,
	})
	if err != nil {
		if s.config.Debug {
			slog.Debug("Failed to encode JSON", "error", err)
		}
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
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

func NewServer(cfg config.Config) (*Server, error) {
	expirationDuration, err := time.ParseDuration(cfg.Expire+"s")
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	client, err := client.New(cfg.HmacKey, cfg.MaxNumber, cfg.Algorithm, cfg.Salt, expirationDuration, cfg.CheckExpire)

	if err != nil {
		return &Server{}, err
	}

	return &Server {
		baseUrl: cfg.BaseUrl,
		port:	cfg.Port,
		client:	*client,
		config: cfg,
	}, nil
}