package handlers_test

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/kvoip/kvoip-pbx/internal/auth"
	"github.com/kvoip/kvoip-pbx/internal/config"
	"github.com/kvoip/kvoip-pbx/internal/dialog"
	"github.com/kvoip/kvoip-pbx/internal/handlers"
	"github.com/kvoip/kvoip-pbx/internal/proxy"
	"github.com/kvoip/kvoip-pbx/internal/session"
	"github.com/kvoip/kvoip-pbx/internal/sip"
)

func newDispatcher() *handlers.Dispatcher {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	cfg := config.Config{
		SIPAdvertisedHost: "127.0.0.1",
		SIPPort:           "5060",
		AuthEnabled:       false,
	}
	return handlers.NewDispatcher(logger, proxy.NewRouter(), session.NewManager(), dialog.NewManager(), cfg, nil, nil)
}

func TestRegisterStoresLocation(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	router := proxy.NewRouter()
	cfg := config.Config{SIPAdvertisedHost: "127.0.0.1", SIPPort: "5060", AuthEnabled: false}
	d := handlers.NewDispatcher(logger, router, session.NewManager(), dialog.NewManager(), cfg, nil, nil)

	var replied []byte
	req := sip.Parse([]byte("REGISTER sip:kvoip.local SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 192.0.2.1:5060;branch=z9hG4bK-1\r\n" +
		"From: <sip:1001@kvoip.local>;tag=abc\r\n" +
		"To: <sip:1001@kvoip.local>\r\n" +
		"Call-ID: call-1@kvoip\r\n" +
		"CSeq: 1 REGISTER\r\n" +
		"Contact: <sip:1001@192.0.2.1:5060>\r\n" +
		"Expires: 3600\r\n" +
		"Content-Length: 0\r\n\r\n"))

	err := d.Handle(handlers.Packet{
		Msg:    req,
		Remote: &net.UDPAddr{IP: net.ParseIP("192.0.2.1"), Port: 5060},
		Reply: func(b []byte) error {
			replied = b
			return nil
		},
		SendTo: func([]byte, *net.UDPAddr) error { return nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(replied), "SIP/2.0 200 OK") {
		t.Fatalf("unexpected response: %s", replied)
	}
	if _, ok := router.Lookup("1001@kvoip.local"); !ok {
		t.Fatal("location não registrada")
	}
}

func TestRegisterDigestChallengeAndSuccess(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	router := proxy.NewRouter()
	cfg := config.Config{
		SIPAdvertisedHost: "127.0.0.1",
		SIPPort:           "5060",
		AuthEnabled:       true,
		AuthRealm:         "kvoip.local",
		SIPUsers:          map[string]string{"1001": "secret"},
	}
	d := handlers.NewDispatcher(logger, router, session.NewManager(), dialog.NewManager(), cfg, nil, nil)

	var replies []string
	base := "REGISTER sip:kvoip.local SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 192.0.2.1:5060;branch=z9hG4bK-1\r\n" +
		"From: <sip:1001@kvoip.local>;tag=abc\r\n" +
		"To: <sip:1001@kvoip.local>\r\n" +
		"Call-ID: call-auth@kvoip\r\n" +
		"CSeq: 1 REGISTER\r\n" +
		"Contact: <sip:1001@192.0.2.1:5060>\r\n" +
		"Expires: 3600\r\n" +
		"Content-Length: 0\r\n\r\n"

	req1 := sip.Parse([]byte(base))
	_ = d.Handle(handlers.Packet{
		Msg:    req1,
		Remote: &net.UDPAddr{IP: net.ParseIP("192.0.2.1"), Port: 5060},
		Reply: func(b []byte) error {
			replies = append(replies, string(b))
			return nil
		},
		SendTo: func([]byte, *net.UDPAddr) error { return nil },
	})
	if len(replies) == 0 || !strings.Contains(replies[0], "401 Unauthorized") {
		t.Fatalf("expected 401, got %#v", replies)
	}
	nonce := extractHeaderValue(replies[0], "WWW-Authenticate")
	nonceVal := extractQuoted(nonce, "nonce")
	if nonceVal == "" {
		t.Fatalf("nonce missing in %s", replies[0])
	}

	params := auth.Params{
		Username: "1001",
		Realm:    "kvoip.local",
		Nonce:    nonceVal,
		URI:      "sip:kvoip.local",
		QOP:      "auth",
		NC:       "00000001",
		CNonce:   "0a4f113b",
	}
	params.Response = auth.ComputeResponse("1001", "kvoip.local", "secret", "REGISTER", params)
	authHeader := fmt.Sprintf(
		`Digest username="1001", realm="kvoip.local", nonce="%s", uri="sip:kvoip.local", response="%s", qop=auth, nc=00000001, cnonce="0a4f113b"`,
		nonceVal,
		params.Response,
	)

	authed := "REGISTER sip:kvoip.local SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 192.0.2.1:5060;branch=z9hG4bK-2\r\n" +
		"From: <sip:1001@kvoip.local>;tag=abc\r\n" +
		"To: <sip:1001@kvoip.local>\r\n" +
		"Call-ID: call-auth@kvoip\r\n" +
		"CSeq: 2 REGISTER\r\n" +
		"Authorization: " + authHeader + "\r\n" +
		"Contact: <sip:1001@192.0.2.1:5060>\r\n" +
		"Expires: 3600\r\n" +
		"Content-Length: 0\r\n\r\n"

	req2 := sip.Parse([]byte(authed))
	var final []byte
	_ = d.Handle(handlers.Packet{
		Msg:    req2,
		Remote: &net.UDPAddr{IP: net.ParseIP("192.0.2.1"), Port: 5060},
		Reply: func(b []byte) error {
			final = b
			return nil
		},
		SendTo: func([]byte, *net.UDPAddr) error { return nil },
	})
	if !strings.Contains(string(final), "200 OK") {
		t.Fatalf("expected 200, got %s", final)
	}
	if _, ok := router.Lookup("1001@kvoip.local"); !ok {
		t.Fatal("not registered after auth")
	}
}

func TestInviteNotFound(t *testing.T) {
	d := newDispatcher()
	var replied []byte
	req := sip.Parse([]byte("INVITE sip:1002@kvoip.local SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 192.0.2.1:5060;branch=z9hG4bK-inv\r\n" +
		"From: <sip:1001@kvoip.local>;tag=from1\r\n" +
		"To: <sip:1002@kvoip.local>\r\n" +
		"Call-ID: invite-404\r\n" +
		"CSeq: 1 INVITE\r\n" +
		"Content-Length: 0\r\n\r\n"))

	err := d.Handle(handlers.Packet{
		Msg:    req,
		Remote: &net.UDPAddr{IP: net.ParseIP("192.0.2.1"), Port: 5060},
		Reply: func(b []byte) error {
			replied = append(replied, b...)
			return nil
		},
		SendTo: func([]byte, *net.UDPAddr) error { return nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(replied), "404 Not Found") {
		t.Fatalf("expected 404, got %s", replied)
	}
}

func TestInviteProxiesToRegisteredContact(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	router := proxy.NewRouter()
	router.Register(proxy.Location{
		AOR:     "1002@kvoip.local",
		Contact: "<sip:1002@198.51.100.2:5070>",
		Expires: 3600,
	})
	cfg := config.Config{SIPAdvertisedHost: "127.0.0.1", SIPPort: "5060", AuthEnabled: false}
	d := handlers.NewDispatcher(logger, router, session.NewManager(), dialog.NewManager(), cfg, nil, nil)

	var replies []string
	var forwarded []byte
	var forwardedTo *net.UDPAddr

	req := sip.Parse([]byte("INVITE sip:1002@kvoip.local SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 192.0.2.1:5060;branch=z9hG4bK-inv\r\n" +
		"From: <sip:1001@kvoip.local>;tag=from1\r\n" +
		"To: <sip:1002@kvoip.local>\r\n" +
		"Call-ID: invite-ok\r\n" +
		"CSeq: 1 INVITE\r\n" +
		"Max-Forwards: 70\r\n" +
		"Content-Type: application/sdp\r\n" +
		"Content-Length: 5\r\n\r\n" +
		"v=0\r\n"))

	err := d.Handle(handlers.Packet{
		Msg:    req,
		Remote: &net.UDPAddr{IP: net.ParseIP("192.0.2.1"), Port: 5060},
		Reply: func(b []byte) error {
			replies = append(replies, string(b))
			return nil
		},
		SendTo: func(b []byte, addr *net.UDPAddr) error {
			forwarded = b
			forwardedTo = addr
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(replies) == 0 || !strings.Contains(replies[0], "100 Trying") {
		t.Fatalf("expected 100 Trying, got %#v", replies)
	}
	if forwardedTo == nil || forwardedTo.Port != 5070 {
		t.Fatalf("forwarded to %#v", forwardedTo)
	}
	body := string(forwarded)
	if !strings.Contains(body, "INVITE sip:1002@198.51.100.2:5070 SIP/2.0") {
		t.Fatalf("bad request-uri: %s", body)
	}
	if !strings.Contains(body, "Via: SIP/2.0/UDP 127.0.0.1:5060;branch=") {
		t.Fatalf("missing top via: %s", body)
	}
}

func extractHeaderValue(msg, name string) string {
	for _, line := range strings.Split(msg, "\r\n") {
		if strings.HasPrefix(strings.ToLower(line), strings.ToLower(name)+":") {
			_, value, _ := strings.Cut(line, ":")
			return strings.TrimSpace(value)
		}
	}
	return ""
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
