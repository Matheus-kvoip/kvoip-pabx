package media

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"
)

// Bridge relays RTP/RTCP between caller and callee legs.
type Bridge struct {
	CallID string

	callerRTP  *net.UDPConn
	callerRTCP *net.UDPConn
	calleeRTP  *net.UDPConn
	calleeRTCP *net.UDPConn

	mu           sync.RWMutex
	callerRemote *net.UDPAddr
	calleeRemote *net.UDPAddr

	logger *slog.Logger
	done   chan struct{}
	once   sync.Once
}

func (b *Bridge) CallerRTPPort() int {
	return udpPort(b.callerRTP)
}

func (b *Bridge) CalleeRTPPort() int {
	return udpPort(b.calleeRTP)
}

func (b *Bridge) SetCallerRemote(ip string, port int) {
	if ip == "" || port <= 0 {
		return
	}
	addr := &net.UDPAddr{IP: net.ParseIP(ip), Port: port}
	if addr.IP == nil {
		return
	}
	b.mu.Lock()
	b.callerRemote = addr
	b.mu.Unlock()
}

func (b *Bridge) SetCalleeRemote(ip string, port int) {
	if ip == "" || port <= 0 {
		return
	}
	addr := &net.UDPAddr{IP: net.ParseIP(ip), Port: port}
	if addr.IP == nil {
		return
	}
	b.mu.Lock()
	b.calleeRemote = addr
	b.mu.Unlock()
}

func (b *Bridge) start() {
	go b.pump(b.callerRTP, true, false)
	go b.pump(b.callerRTCP, true, true)
	go b.pump(b.calleeRTP, false, false)
	go b.pump(b.calleeRTCP, false, true)
}

func (b *Bridge) pump(conn *net.UDPConn, fromCaller, rtcp bool) {
	if conn == nil {
		return
	}
	buf := make([]byte, 2048)
	for {
		_ = conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		n, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			select {
			case <-b.done:
				return
			default:
				// idle timeout — keep waiting until closed
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					continue
				}
				return
			}
		}

		b.learnRemote(fromCaller, rtcp, remote)

		b.mu.RLock()
		var dest *net.UDPAddr
		var out *net.UDPConn
		if fromCaller {
			dest = b.calleeRemote
			if rtcp {
				out = b.calleeRTCP
				if dest != nil {
					dest = &net.UDPAddr{IP: dest.IP, Port: dest.Port + 1}
				}
			} else {
				out = b.calleeRTP
			}
		} else {
			dest = b.callerRemote
			if rtcp {
				out = b.callerRTCP
				if dest != nil {
					dest = &net.UDPAddr{IP: dest.IP, Port: dest.Port + 1}
				}
			} else {
				out = b.callerRTP
			}
		}
		b.mu.RUnlock()

		if dest == nil || out == nil {
			continue
		}
		if _, err := out.WriteToUDP(buf[:n], dest); err != nil {
			b.logger.Debug("rtp write failed", "call_id", b.CallID, "err", err)
		}
	}
}

func (b *Bridge) learnRemote(fromCaller, rtcp bool, remote *net.UDPAddr) {
	if remote == nil || rtcp {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if fromCaller {
		if b.callerRemote == nil || b.callerRemote.Port != remote.Port || !b.callerRemote.IP.Equal(remote.IP) {
			b.callerRemote = cloneAddr(remote)
			b.logger.Debug("rtp caller remote", "call_id", b.CallID, "addr", remote.String())
		}
	} else if b.calleeRemote == nil || b.calleeRemote.Port != remote.Port || !b.calleeRemote.IP.Equal(remote.IP) {
		b.calleeRemote = cloneAddr(remote)
		b.logger.Debug("rtp callee remote", "call_id", b.CallID, "addr", remote.String())
	}
}

func (b *Bridge) Close() {
	b.once.Do(func() {
		close(b.done)
		_ = b.callerRTP.Close()
		_ = b.callerRTCP.Close()
		_ = b.calleeRTP.Close()
		_ = b.calleeRTCP.Close()
	})
}

func udpPort(conn *net.UDPConn) int {
	if conn == nil {
		return 0
	}
	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok || addr == nil {
		return 0
	}
	return addr.Port
}

func cloneAddr(addr *net.UDPAddr) *net.UDPAddr {
	if addr == nil {
		return nil
	}
	ip := make(net.IP, len(addr.IP))
	copy(ip, addr.IP)
	return &net.UDPAddr{IP: ip, Port: addr.Port, Zone: addr.Zone}
}

func bindPair(host string, port int) (rtp, rtcp *net.UDPConn, err error) {
	ip := net.ParseIP(host)
	if host == "0.0.0.0" || host == "" {
		ip = net.IPv4zero
	}
	rtp, err = net.ListenUDP("udp", &net.UDPAddr{IP: ip, Port: port})
	if err != nil {
		return nil, nil, err
	}
	rtcp, err = net.ListenUDP("udp", &net.UDPAddr{IP: ip, Port: port + 1})
	if err != nil {
		_ = rtp.Close()
		return nil, nil, err
	}
	return rtp, rtcp, nil
}

func closeQuiet(conns ...*net.UDPConn) {
	for _, c := range conns {
		if c != nil {
			_ = c.Close()
		}
	}
}

func fmtErr(op string, err error) error {
	return fmt.Errorf("media %s: %w", op, err)
}
