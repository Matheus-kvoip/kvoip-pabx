package media_test

import (
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/kvoip/kvoip-pbx/internal/media"
)

func TestBridgeRelaysRTP(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	engine := media.NewEngine(logger, "127.0.0.1", "127.0.0.1", 12000, 12100)
	bridge, err := engine.Open("call-1")
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close("call-1")

	callerConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0})
	if err != nil {
		t.Fatal(err)
	}
	defer callerConn.Close()

	calleeConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0})
	if err != nil {
		t.Fatal(err)
	}
	defer calleeConn.Close()

	callerLocal := callerConn.LocalAddr().(*net.UDPAddr)
	calleeLocal := calleeConn.LocalAddr().(*net.UDPAddr)
	bridge.SetCallerRemote("127.0.0.1", callerLocal.Port)
	bridge.SetCalleeRemote("127.0.0.1", calleeLocal.Port)

	payload := []byte("hello-rtp")
	dst := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: bridge.CallerRTPPort()}
	if _, err := callerConn.WriteToUDP(payload, dst); err != nil {
		t.Fatal(err)
	}

	_ = calleeConn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 64)
	n, _, err := calleeConn.ReadFromUDP(buf)
	if err != nil {
		t.Fatalf("callee did not receive relayed packet: %v", err)
	}
	if string(buf[:n]) != string(payload) {
		t.Fatalf("got %q", buf[:n])
	}
}
