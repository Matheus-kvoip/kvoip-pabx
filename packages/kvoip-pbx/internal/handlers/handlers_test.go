package handlers_test

import (
	"log/slog"
	"net"
	"os"
	"strings"
	"testing"

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
	}
	return handlers.NewDispatcher(logger, proxy.NewRouter(), session.NewManager(), dialog.NewManager(), cfg)
}

func TestRegisterStoresLocation(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	router := proxy.NewRouter()
	cfg := config.Config{SIPAdvertisedHost: "127.0.0.1", SIPPort: "5060"}
	d := handlers.NewDispatcher(logger, router, session.NewManager(), dialog.NewManager(), cfg)

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
	cfg := config.Config{SIPAdvertisedHost: "127.0.0.1", SIPPort: "5060"}
	d := handlers.NewDispatcher(logger, router, session.NewManager(), dialog.NewManager(), cfg)

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
