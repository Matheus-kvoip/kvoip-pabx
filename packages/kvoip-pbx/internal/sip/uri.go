package sip

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ExtractURI returns the SIP URI inside angle brackets, or the raw token.
func ExtractURI(value string) string {
	value = strings.TrimSpace(value)
	if start := strings.Index(value, "<"); start >= 0 {
		if end := strings.Index(value[start:], ">"); end > 0 {
			return value[start+1 : start+end]
		}
	}
	if i := strings.IndexAny(value, " ;"); i >= 0 {
		return value[:i]
	}
	return value
}

// ExtractAOR returns user@host from a SIP URI / From / To header.
func ExtractAOR(value string) string {
	uri := ExtractURI(value)
	uri = strings.TrimPrefix(uri, "sip:")
	uri = strings.TrimPrefix(uri, "sips:")
	if i := strings.Index(uri, ";"); i >= 0 {
		uri = uri[:i]
	}
	return uri
}

// ExtractUser returns the user part of an AOR/URI.
func ExtractUser(value string) string {
	aor := ExtractAOR(value)
	if i := strings.Index(aor, "@"); i >= 0 {
		return aor[:i]
	}
	return aor
}

// ContactExpires reads expires from Contact params or Expires header.
func ContactExpires(contact string, expiresHeader string, fallback int) int {
	lower := strings.ToLower(contact)
	if i := strings.Index(lower, "expires="); i >= 0 {
		raw := contact[i+len("expires="):]
		if end := strings.IndexAny(raw, ";, \t"); end >= 0 {
			raw = raw[:end]
		}
		if n, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil {
			return n
		}
	}
	if expiresHeader != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(expiresHeader)); err == nil {
			return n
		}
	}
	return fallback
}

// UDPAddrFromURI resolves host:port from a SIP URI / Contact.
func UDPAddrFromURI(value string) (*net.UDPAddr, error) {
	uri := ExtractURI(value)
	uri = strings.TrimPrefix(uri, "sip:")
	uri = strings.TrimPrefix(uri, "sips:")
	if i := strings.Index(uri, ";"); i >= 0 {
		uri = uri[:i]
	}
	if at := strings.LastIndex(uri, "@"); at >= 0 {
		uri = uri[at+1:]
	}
	if !strings.Contains(uri, ":") {
		uri = net.JoinHostPort(uri, "5060")
	}
	return net.ResolveUDPAddr("udp", uri)
}

// NewBranch creates a Via branch parameter.
func NewBranch(seed string) string {
	seed = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, seed)
	if len(seed) > 12 {
		seed = seed[:12]
	}
	if seed == "" {
		seed = "kvoip"
	}
	return fmt.Sprintf("z9hG4bK-%s", seed)
}
