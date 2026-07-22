package sip_test

import (
	"strings"
	"testing"

	"github.com/kvoip/kvoip-pbx/internal/sip"
)

func TestParseRegister(t *testing.T) {
	raw := []byte("REGISTER sip:kvoip.local SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 192.0.2.1:5060;branch=z9hG4bK-1\r\n" +
		"Via: SIP/2.0/UDP 198.51.100.1:5060;branch=z9hG4bK-2\r\n" +
		"From: <sip:1001@kvoip.local>;tag=abc\r\n" +
		"To: <sip:1001@kvoip.local>\r\n" +
		"Call-ID: call-1@kvoip\r\n" +
		"CSeq: 1 REGISTER\r\n" +
		"Contact: <sip:1001@192.0.2.1:5060>\r\n" +
		"Expires: 3600\r\n" +
		"Content-Length: 0\r\n\r\n")

	msg := sip.Parse(raw)
	if !msg.IsRequest {
		t.Fatalf("esperado request")
	}
	if msg.Method != sip.MethodRegister {
		t.Fatalf("method=%q", msg.Method)
	}
	if msg.RequestURI != "sip:kvoip.local" {
		t.Fatalf("uri=%q", msg.RequestURI)
	}
	if len(msg.Vias) != 2 {
		t.Fatalf("vias=%d", len(msg.Vias))
	}
}

func TestBuildResponseCopiesDialogHeaders(t *testing.T) {
	req := sip.Parse([]byte("OPTIONS sip:kvoip.local SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 192.0.2.1:5060;branch=z9hG4bK-opt\r\n" +
		"From: <sip:scanner@example.com>;tag=1\r\n" +
		"To: <sip:kvoip.local>\r\n" +
		"Call-ID: opt-1\r\n" +
		"CSeq: 1 OPTIONS\r\n" +
		"Content-Length: 0\r\n\r\n"))

	res := string(sip.BuildResponse(req, 200, "OK", map[string]string{
		"Allow": "OPTIONS, REGISTER",
	}))

	for _, want := range []string{
		"SIP/2.0 200 OK",
		"Via: SIP/2.0/UDP 192.0.2.1:5060;branch=z9hG4bK-opt",
		"From: <sip:scanner@example.com>;tag=1",
		"To: <sip:kvoip.local>;tag=kvoip-",
		"Call-ID: opt-1",
		"CSeq: 1 OPTIONS",
		"Allow: OPTIONS, REGISTER",
	} {
		if !strings.Contains(res, want) {
			t.Fatalf("resposta sem %q\n%s", want, res)
		}
	}
}

func TestExtractAORAndUDPAddr(t *testing.T) {
	aor := sip.ExtractAOR(`"Ana" <sip:1001@kvoip.local>;tag=x`)
	if aor != "1001@kvoip.local" {
		t.Fatalf("aor=%q", aor)
	}
	addr, err := sip.UDPAddrFromURI(`<sip:1002@198.51.100.2:5070>`)
	if err != nil {
		t.Fatal(err)
	}
	if addr.Port != 5070 || addr.IP.String() != "198.51.100.2" {
		t.Fatalf("addr=%v", addr)
	}
}

func TestForwardResponseStripsTopVia(t *testing.T) {
	res := sip.Parse([]byte("SIP/2.0 180 Ringing\r\n" +
		"Via: SIP/2.0/UDP 127.0.0.1:5060;branch=z9hG4bK-pbx\r\n" +
		"Via: SIP/2.0/UDP 192.0.2.1:5060;branch=z9hG4bK-ua\r\n" +
		"From: <sip:1001@kvoip.local>;tag=from1\r\n" +
		"To: <sip:1002@kvoip.local>;tag=to1\r\n" +
		"Call-ID: c1\r\n" +
		"CSeq: 1 INVITE\r\n" +
		"Content-Length: 0\r\n\r\n"))
	out := string(sip.ForwardResponse(res))
	if strings.Contains(out, "127.0.0.1:5060") {
		t.Fatalf("top via should be stripped: %s", out)
	}
	if !strings.Contains(out, "192.0.2.1:5060") {
		t.Fatalf("caller via missing: %s", out)
	}
}
