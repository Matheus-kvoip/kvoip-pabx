package session

import (
	"sync"
	"time"
)

// State is the call media/session state.
type State string

const (
	StateIdle     State = "idle"
	StateRinging  State = "ringing"
	StateAnswered State = "answered"
	StateHeld     State = "held"
	StateEnded    State = "ended"
)

// Call is an active or recent call session in the PBX.
type Call struct {
	ID         string     `json:"id"`
	From       string     `json:"from"`
	To         string     `json:"to"`
	State      State      `json:"state"`
	StartedAt  time.Time  `json:"startedAt"`
	AnsweredAt *time.Time `json:"answeredAt,omitempty"`
	EndedAt    *time.Time `json:"endedAt,omitempty"`
}

// Manager keeps in-memory call sessions.
type Manager struct {
	mu    sync.RWMutex
	calls map[string]*Call
}

func NewManager() *Manager {
	return &Manager{calls: make(map[string]*Call)}
}

func (m *Manager) Upsert(call *Call) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if existing, ok := m.calls[call.ID]; ok {
		if call.StartedAt.IsZero() {
			call.StartedAt = existing.StartedAt
		}
		if call.AnsweredAt == nil {
			call.AnsweredAt = existing.AnsweredAt
		}
		if call.EndedAt == nil {
			call.EndedAt = existing.EndedAt
		}
	}
	if call.StartedAt.IsZero() {
		call.StartedAt = time.Now().UTC()
	}
	copied := *call
	m.calls[call.ID] = &copied
}

func (m *Manager) Get(id string) (*Call, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	call, ok := m.calls[id]
	if !ok {
		return nil, false
	}
	copied := *call
	return &copied, true
}

func (m *Manager) List() []Call {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Call, 0, len(m.calls))
	for _, call := range m.calls {
		out = append(out, *call)
	}
	return out
}

func (m *Manager) ActiveCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	count := 0
	for _, call := range m.calls {
		if call.State != StateEnded {
			count++
		}
	}
	return count
}

func (m *Manager) MarkAnswered(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if call, ok := m.calls[id]; ok {
		call.State = StateAnswered
		now := time.Now().UTC()
		call.AnsweredAt = &now
	}
}

func (m *Manager) MarkEnded(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if call, ok := m.calls[id]; ok {
		call.State = StateEnded
		now := time.Now().UTC()
		call.EndedAt = &now
	}
}
