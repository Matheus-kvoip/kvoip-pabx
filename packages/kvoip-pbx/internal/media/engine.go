package media

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
)

// Engine allocates RTP bridges for active calls.
type Engine struct {
	logger        *slog.Logger
	advertiseHost string
	bindHost      string
	portMin       int
	portMax       int

	mu      sync.Mutex
	next    int
	bridges map[string]*Bridge
}

func NewEngine(logger *slog.Logger, advertiseHost, bindHost string, portMin, portMax int) *Engine {
	if portMin <= 0 {
		portMin = 10000
	}
	if portMax <= portMin+4 {
		portMax = portMin + 4000
	}
	if bindHost == "" {
		bindHost = "0.0.0.0"
	}
	if advertiseHost == "" {
		advertiseHost = "127.0.0.1"
	}
	if portMin%2 != 0 {
		portMin++
	}
	return &Engine{
		logger:        logger,
		advertiseHost: advertiseHost,
		bindHost:      bindHost,
		portMin:       portMin,
		portMax:       portMax,
		next:          portMin,
		bridges:       make(map[string]*Bridge),
	}
}

func (e *Engine) AdvertiseHost() string {
	return e.advertiseHost
}

// Open creates (or replaces) a bridge for callID.
func (e *Engine) Open(callID string) (*Bridge, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if old, ok := e.bridges[callID]; ok {
		old.Close()
		delete(e.bridges, callID)
	}

	callerRTP, callerRTCP, err := e.bindNextPair()
	if err != nil {
		return nil, fmtErr("caller bind", err)
	}
	calleeRTP, calleeRTCP, err := e.bindNextPair()
	if err != nil {
		closeQuiet(callerRTP, callerRTCP)
		return nil, fmtErr("callee bind", err)
	}

	b := &Bridge{
		CallID:     callID,
		callerRTP:  callerRTP,
		callerRTCP: callerRTCP,
		calleeRTP:  calleeRTP,
		calleeRTCP: calleeRTCP,
		logger:     e.logger,
		done:       make(chan struct{}),
	}
	b.start()
	e.bridges[callID] = b
	e.logger.Info("rtp bridge open",
		"call_id", callID,
		"caller_port", b.CallerRTPPort(),
		"callee_port", b.CalleeRTPPort(),
	)
	return b, nil
}

func (e *Engine) Get(callID string) (*Bridge, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	b, ok := e.bridges[callID]
	return b, ok
}

func (e *Engine) Close(callID string) {
	e.mu.Lock()
	b, ok := e.bridges[callID]
	if ok {
		delete(e.bridges, callID)
	}
	e.mu.Unlock()
	if ok {
		b.Close()
		e.logger.Info("rtp bridge close", "call_id", callID)
	}
}

func (e *Engine) Active() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.bridges)
}

func (e *Engine) bindNextPair() (*net.UDPConn, *net.UDPConn, error) {
	for attempt := 0; attempt < ((e.portMax-e.portMin)/2)+2; attempt++ {
		port := e.next
		e.next += 2
		if e.next+1 > e.portMax {
			e.next = e.portMin
		}
		rtp, rtcp, err := bindPair(e.bindHost, port)
		if err == nil {
			return rtp, rtcp, nil
		}
	}
	return nil, nil, fmt.Errorf("no free RTP ports in %d-%d", e.portMin, e.portMax)
}
