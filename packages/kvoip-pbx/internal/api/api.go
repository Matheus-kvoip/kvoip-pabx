package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/kvoip/kvoip-pbx/internal/config"
	"github.com/kvoip/kvoip-pbx/internal/proxy"
	"github.com/kvoip/kvoip-pbx/internal/session"
	"github.com/kvoip/kvoip-pbx/pkg/version"
)

// Server exposes a small HTTP control-plane API for Nest/front.
type Server struct {
	cfg      config.Config
	logger   *slog.Logger
	router   *proxy.Router
	sessions *session.Manager
	started  time.Time
}

func New(cfg config.Config, logger *slog.Logger, router *proxy.Router, sessions *session.Manager) *Server {
	return &Server{
		cfg:      cfg,
		logger:   logger,
		router:   router,
		sessions: sessions,
		started:  time.Now().UTC(),
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/v1/registrations", s.handleRegistrations)
	mux.HandleFunc("/v1/calls", s.handleCalls)
	return withCORS(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "ok",
		"service":   s.cfg.ServiceName,
		"version":   version.Version,
		"uptimeSec": int(time.Since(s.started).Seconds()),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"sip":       s.cfg.ListenAddr(),
		"bindings":  s.router.Count(),
		"activeCalls": s.sessions.ActiveCount(),
	})
}

type registrationDTO struct {
	AOR       string `json:"aor"`
	Number    string `json:"number"`
	Contact   string `json:"contact"`
	Expires   int    `json:"expires"`
	UpdatedAt string `json:"updatedAt"`
}

func (s *Server) handleRegistrations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	items := s.router.List()
	out := make([]registrationDTO, 0, len(items))
	for _, loc := range items {
		number := loc.AOR
		if i := strings.Index(number, "@"); i >= 0 {
			number = number[:i]
		}
		out = append(out, registrationDTO{
			AOR:       loc.AOR,
			Number:    number,
			Contact:   loc.Contact,
			Expires:   loc.Expires,
			UpdatedAt: loc.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}
	writeJSON(w, http.StatusOK, out)
}

type callDTO struct {
	ID          string `json:"id"`
	Direction   string `json:"direction"`
	State       string `json:"state"`
	From        string `json:"from"`
	To          string `json:"to"`
	StartedAt   string `json:"startedAt"`
	AnsweredAt  string `json:"answeredAt,omitempty"`
	EndedAt     string `json:"endedAt,omitempty"`
	DurationSec int    `json:"durationSec"`
}

func (s *Server) handleCalls(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	activeOnly := r.URL.Query().Get("active") == "true"
	items := s.sessions.List()
	out := make([]callDTO, 0, len(items))
	now := time.Now().UTC()
	for _, call := range items {
		if activeOnly && call.State == session.StateEnded {
			continue
		}
		fromUser := trimUser(call.From)
		toUser := trimUser(call.To)
		end := now
		if call.EndedAt != nil {
			end = *call.EndedAt
		}
		duration := int(end.Sub(call.StartedAt).Seconds())
		if duration < 0 {
			duration = 0
		}
		dto := callDTO{
			ID:          call.ID,
			Direction:   "internal",
			State:       string(call.State),
			From:        fromUser,
			To:          toUser,
			StartedAt:   call.StartedAt.UTC().Format(time.RFC3339),
			DurationSec: duration,
		}
		if call.AnsweredAt != nil {
			dto.AnsweredAt = call.AnsweredAt.UTC().Format(time.RFC3339)
		}
		if call.EndedAt != nil {
			dto.EndedAt = call.EndedAt.UTC().Format(time.RFC3339)
		}
		out = append(out, dto)
	}
	writeJSON(w, http.StatusOK, out)
}

func trimUser(aor string) string {
	aor = strings.TrimPrefix(aor, "sip:")
	if i := strings.Index(aor, "@"); i >= 0 {
		return aor[:i]
	}
	return aor
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
