package proxy

import (
	"strings"
	"sync"
	"time"
)

// Location is a registered SIP contact for an extension/AOR.
type Location struct {
	AOR       string
	Contact   string
	Expires   int
	UpdatedAt time.Time
}

// Router selects the next hop for a request and stores registrations.
type Router struct {
	mu        sync.RWMutex
	locations map[string]Location
}

func NewRouter() *Router {
	return &Router{locations: make(map[string]Location)}
}

func (r *Router) Register(loc Location) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if loc.UpdatedAt.IsZero() {
		loc.UpdatedAt = time.Now().UTC()
	}
	if loc.Expires <= 0 {
		delete(r.locations, loc.AOR)
		return
	}
	r.locations[loc.AOR] = loc
}

func (r *Router) Unregister(aor string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.locations, aor)
}

func (r *Router) Lookup(aor string) (Location, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	loc, ok := r.locations[aor]
	return loc, ok
}

// LookupFlexible matches full AOR or user part (1001 / 1001@domain).
func (r *Router) LookupFlexible(aorOrURI string) (Location, bool) {
	aor := aorOrURI
	if strings.Contains(aor, ":") && strings.Contains(aor, "@") {
		// likely sip:user@host
		aor = strings.TrimPrefix(aor, "sip:")
		aor = strings.TrimPrefix(aor, "sips:")
		if i := strings.Index(aor, ";"); i >= 0 {
			aor = aor[:i]
		}
	}
	if loc, ok := r.Lookup(aor); ok {
		return loc, true
	}

	user := aor
	if i := strings.Index(user, "@"); i >= 0 {
		user = user[:i]
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for key, loc := range r.locations {
		u := key
		if i := strings.Index(u, "@"); i >= 0 {
			u = u[:i]
		}
		if u == user {
			return loc, true
		}
	}
	return Location{}, false
}

func (r *Router) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.locations)
}

// List returns a snapshot of all registrations.
func (r *Router) List() []Location {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Location, 0, len(r.locations))
	for _, loc := range r.locations {
		out = append(out, loc)
	}
	return out
}
