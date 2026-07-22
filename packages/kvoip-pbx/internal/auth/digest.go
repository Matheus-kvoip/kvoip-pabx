package auth

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Credentials maps SIP username -> plain password (MVP).
type Credentials map[string]string

// Digest handles SIP Digest challenges and validation.
type Digest struct {
	Realm string
	Users Credentials

	mu      sync.Mutex
	nonces  map[string]time.Time
	ttl     time.Duration
}

func NewDigest(realm string, users Credentials) *Digest {
	if realm == "" {
		realm = "kvoip.local"
	}
	return &Digest{
		Realm:  realm,
		Users:  users,
		nonces: make(map[string]time.Time),
		ttl:    5 * time.Minute,
	}
}

func (d *Digest) ReplaceUsers(users Credentials) {
	d.mu.Lock()
	defer d.mu.Unlock()
	cp := Credentials{}
	for k, v := range users {
		cp[k] = v
	}
	d.Users = cp
}

func (d *Digest) UserCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.Users)
}

// ChallengeHeader builds a WWW-Authenticate value.
func (d *Digest) ChallengeHeader() string {
	nonce := d.issueNonce()
	return fmt.Sprintf(
		`Digest realm="%s", nonce="%s", algorithm=MD5, qop="auth"`,
		d.Realm,
		nonce,
	)
}

func (d *Digest) issueNonce() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	nonce := hex.EncodeToString(buf)
	d.mu.Lock()
	d.nonces[nonce] = time.Now().UTC().Add(d.ttl)
	d.cleanupLocked(time.Now().UTC())
	d.mu.Unlock()
	return nonce
}

func (d *Digest) cleanupLocked(now time.Time) {
	for n, exp := range d.nonces {
		if now.After(exp) {
			delete(d.nonces, n)
		}
	}
}

func (d *Digest) nonceValid(nonce string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := time.Now().UTC()
	d.cleanupLocked(now)
	exp, ok := d.nonces[nonce]
	return ok && !now.After(exp)
}

// Params is a parsed Authorization: Digest ... header.
type Params struct {
	Username string
	Realm    string
	Nonce    string
	URI      string
	Response string
	QOP      string
	NC       string
	CNonce   string
	Algorithm string
}

// ParseAuthorization parses `Authorization: Digest ...` or raw Digest value.
func ParseAuthorization(header string) (Params, bool) {
	header = strings.TrimSpace(header)
	if header == "" {
		return Params{}, false
	}
	lower := strings.ToLower(header)
	if strings.HasPrefix(lower, "digest ") {
		header = strings.TrimSpace(header[len("Digest "):])
	} else if !strings.Contains(lower, "username=") {
		return Params{}, false
	}

	out := Params{}
	for _, part := range splitAuthParts(header) {
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		key = strings.ToLower(strings.TrimSpace(key))
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"`)
		switch key {
		case "username":
			out.Username = value
		case "realm":
			out.Realm = value
		case "nonce":
			out.Nonce = value
		case "uri":
			out.URI = value
		case "response":
			out.Response = value
		case "qop":
			out.QOP = value
		case "nc":
			out.NC = value
		case "cnonce":
			out.CNonce = value
		case "algorithm":
			out.Algorithm = value
		}
	}
	if out.Username == "" || out.Nonce == "" || out.Response == "" || out.URI == "" {
		return Params{}, false
	}
	return out, true
}

func splitAuthParts(header string) []string {
	var parts []string
	var b strings.Builder
	inQuotes := false
	for i := 0; i < len(header); i++ {
		ch := header[i]
		switch ch {
		case '"':
			inQuotes = !inQuotes
			b.WriteByte(ch)
		case ',':
			if inQuotes {
				b.WriteByte(ch)
			} else {
				parts = append(parts, strings.TrimSpace(b.String()))
				b.Reset()
			}
		default:
			b.WriteByte(ch)
		}
	}
	if s := strings.TrimSpace(b.String()); s != "" {
		parts = append(parts, s)
	}
	return parts
}

// Validate checks Digest credentials for a SIP method.
func (d *Digest) Validate(method, authHeader string) (username string, ok bool) {
	params, parsed := ParseAuthorization(authHeader)
	if !parsed {
		return "", false
	}
	if params.Realm != "" && params.Realm != d.Realm {
		return "", false
	}
	if !d.nonceValid(params.Nonce) {
		return "", false
	}
	d.mu.Lock()
	password, exists := d.Users[params.Username]
	d.mu.Unlock()
	if !exists {
		return "", false
	}

	expected := ComputeResponse(params.Username, d.Realm, password, method, params)
	if !strings.EqualFold(expected, params.Response) {
		return "", false
	}
	return params.Username, true
}

// ComputeResponse builds the Digest response (MD5).
func ComputeResponse(username, realm, password, method string, p Params) string {
	ha1 := md5Hex(fmt.Sprintf("%s:%s:%s", username, realm, password))
	ha2 := md5Hex(fmt.Sprintf("%s:%s", method, p.URI))
	if strings.EqualFold(p.QOP, "auth") {
		return md5Hex(fmt.Sprintf("%s:%s:%s:%s:%s:%s", ha1, p.Nonce, p.NC, p.CNonce, p.QOP, ha2))
	}
	return md5Hex(fmt.Sprintf("%s:%s:%s", ha1, p.Nonce, ha2))
}

func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// ParseUsers parses "1001:pass,1002:pass" into Credentials.
func ParseUsers(raw string) Credentials {
	out := Credentials{}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return out
	}
	for _, item := range strings.Split(raw, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		user, pass, ok := strings.Cut(item, ":")
		if !ok {
			continue
		}
		user = strings.TrimSpace(user)
		pass = strings.TrimSpace(pass)
		if user == "" {
			continue
		}
		out[user] = pass
	}
	return out
}
