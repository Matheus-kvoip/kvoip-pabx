package auth_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kvoip/kvoip-pbx/internal/auth"
)

func TestDigestChallengeAndValidate(t *testing.T) {
	d := auth.NewDigest("kvoip.local", auth.Credentials{"1001": "secret"})
	challenge := d.ChallengeHeader()
	if !strings.Contains(challenge, `realm="kvoip.local"`) {
		t.Fatalf("challenge=%s", challenge)
	}
	nonce := extractQuoted(challenge, "nonce")
	if nonce == "" {
		t.Fatal("nonce missing")
	}

	params := auth.Params{
		Username: "1001",
		Realm:    "kvoip.local",
		Nonce:    nonce,
		URI:      "sip:kvoip.local",
		QOP:      "auth",
		NC:       "00000001",
		CNonce:   "xyz",
	}
	params.Response = auth.ComputeResponse("1001", "kvoip.local", "secret", "REGISTER", params)

	header := fmt.Sprintf(
		`Digest username="1001", realm="kvoip.local", nonce="%s", uri="sip:kvoip.local", response="%s", qop=auth, nc=00000001, cnonce="xyz"`,
		nonce,
		params.Response,
	)
	user, ok := d.Validate("REGISTER", header)
	if !ok || user != "1001" {
		t.Fatalf("validate failed user=%q ok=%v", user, ok)
	}
}

func TestParseUsers(t *testing.T) {
	users := auth.ParseUsers("1001:a, 1002:b")
	if users["1001"] != "a" || users["1002"] != "b" {
		t.Fatalf("%v", users)
	}
}

func extractQuoted(header, key string) string {
	needle := key + `="`
	i := strings.Index(header, needle)
	if i < 0 {
		return ""
	}
	rest := header[i+len(needle):]
	j := strings.Index(rest, `"`)
	if j < 0 {
		return ""
	}
	return rest[:j]
}
