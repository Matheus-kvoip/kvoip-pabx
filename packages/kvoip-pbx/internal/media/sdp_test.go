package media_test

import (
	"strings"
	"testing"

	"github.com/kvoip/kvoip-pbx/internal/media"
)

func TestParseAndRewriteAudio(t *testing.T) {
	sdp := "v=0\r\n" +
		"o=- 0 0 IN IP4 192.0.2.10\r\n" +
		"s=-\r\n" +
		"c=IN IP4 192.0.2.10\r\n" +
		"t=0 0\r\n" +
		"m=audio 4000 RTP/AVP 0 8\r\n" +
		"a=rtpmap:0 PCMU/8000\r\n"

	audio, ok := media.ParseAudio(sdp)
	if !ok {
		t.Fatal("parse failed")
	}
	if audio.IP != "192.0.2.10" || audio.Port != 4000 {
		t.Fatalf("unexpected audio: %#v", audio)
	}

	out := media.RewriteAudio(sdp, "203.0.113.5", 10000)
	if !strings.Contains(out, "c=IN IP4 203.0.113.5") {
		t.Fatalf("missing rewritten c=: %s", out)
	}
	if !strings.Contains(out, "m=audio 10000 RTP/AVP 0 8") {
		t.Fatalf("missing rewritten m=: %s", out)
	}
	parsed, ok := media.ParseAudio(out)
	if !ok || parsed.IP != "203.0.113.5" || parsed.Port != 10000 {
		t.Fatalf("rewritten parse: %#v ok=%v", parsed, ok)
	}
}
