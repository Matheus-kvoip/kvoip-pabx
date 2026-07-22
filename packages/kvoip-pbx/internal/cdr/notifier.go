package cdr

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/kvoip/kvoip-pbx/internal/session"
)

// Notifier posts ended calls to the Nest CDR webhook.
type Notifier struct {
	url    string
	secret string
	logger *slog.Logger
	client *http.Client
}

func NewNotifier(url, secret string, logger *slog.Logger) *Notifier {
	url = strings.TrimSpace(url)
	if url == "" {
		return nil
	}
	return &Notifier{
		url:    strings.TrimRight(url, "/"),
		secret: secret,
		logger: logger,
		client: &http.Client{Timeout: 3 * time.Second},
	}
}

type payload struct {
	ID          string  `json:"id"`
	Direction   string  `json:"direction"`
	State       string  `json:"state"`
	From        string  `json:"from"`
	To          string  `json:"to"`
	StartedAt   string  `json:"startedAt"`
	AnsweredAt  *string `json:"answeredAt,omitempty"`
	EndedAt     *string `json:"endedAt,omitempty"`
	DurationSec int     `json:"durationSec"`
}

func (n *Notifier) NotifyEnded(call session.Call) {
	if n == nil {
		return
	}
	from := trimUser(call.From)
	to := trimUser(call.To)
	end := time.Now().UTC()
	if call.EndedAt != nil {
		end = *call.EndedAt
	}
	duration := int(end.Sub(call.StartedAt).Seconds())
	if duration < 0 {
		duration = 0
	}
	body := payload{
		ID:          call.ID,
		Direction:   "internal",
		State:       string(session.StateEnded),
		From:        from,
		To:          to,
		StartedAt:   call.StartedAt.UTC().Format(time.RFC3339),
		DurationSec: duration,
	}
	if call.AnsweredAt != nil {
		s := call.AnsweredAt.UTC().Format(time.RFC3339)
		body.AnsweredAt = &s
	}
	ended := end.Format(time.RFC3339)
	body.EndedAt = &ended

	raw, err := json.Marshal(body)
	if err != nil {
		return
	}

	go func() {
		req, err := http.NewRequest(http.MethodPost, n.url, bytes.NewReader(raw))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		if n.secret != "" {
			req.Header.Set("X-CDR-Secret", n.secret)
		}
		res, err := n.client.Do(req)
		if err != nil {
			n.logger.Warn("cdr webhook falhou", "err", err, "call_id", call.ID)
			return
		}
		defer res.Body.Close()
		if res.StatusCode >= 300 {
			n.logger.Warn("cdr webhook status", "status", res.StatusCode, "call_id", call.ID)
			return
		}
		n.logger.Info("cdr enviado", "call_id", call.ID)
	}()
}

func trimUser(aor string) string {
	aor = strings.TrimPrefix(aor, "sip:")
	if i := strings.Index(aor, "@"); i >= 0 {
		return aor[:i]
	}
	return aor
}
