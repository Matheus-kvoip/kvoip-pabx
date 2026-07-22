package dialog

import (
	"net"
	"sync"
)

// State represents a SIP dialog lifecycle state (RFC 3261).
type State string

const (
	StateEarly      State = "early"
	StateConfirmed  State = "confirmed"
	StateTerminated State = "terminated"
)

// Leg tracks a proxied call between caller and callee.
type Leg struct {
	CallID      string
	State       State
	CallerAddr  *net.UDPAddr
	CalleeAddr  *net.UDPAddr
	CalleeURI   string
	From        string
	To          string
	Branch      string
}

// Manager stores active dialogs by Call-ID.
type Manager struct {
	mu   sync.RWMutex
	legs map[string]*Leg
}

func NewManager() *Manager {
	return &Manager{legs: make(map[string]*Leg)}
}

func (m *Manager) Put(leg *Leg) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.legs[leg.CallID] = leg
}

func (m *Manager) Get(callID string) (*Leg, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	leg, ok := m.legs[callID]
	return leg, ok
}

func (m *Manager) Delete(callID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.legs, callID)
}
