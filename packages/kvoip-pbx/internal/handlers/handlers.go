package handlers

import (
	"fmt"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/kvoip/kvoip-pbx/internal/auth"
	"github.com/kvoip/kvoip-pbx/internal/cdr"
	"github.com/kvoip/kvoip-pbx/internal/config"
	"github.com/kvoip/kvoip-pbx/internal/dialog"
	"github.com/kvoip/kvoip-pbx/internal/media"
	"github.com/kvoip/kvoip-pbx/internal/proxy"
	"github.com/kvoip/kvoip-pbx/internal/session"
	"github.com/kvoip/kvoip-pbx/internal/sip"
)

// Packet is one inbound SIP datagram with transport helpers.
type Packet struct {
	Msg    sip.Message
	Remote *net.UDPAddr
	Reply  func([]byte) error
	SendTo func([]byte, *net.UDPAddr) error
}

// Dispatcher routes SIP messages (requests and responses).
type Dispatcher struct {
	logger   *slog.Logger
	router   *proxy.Router
	sessions *session.Manager
	dialogs  *dialog.Manager
	cfg      config.Config
	digest   *auth.Digest
	media    *media.Engine
	cdr      *cdr.Notifier
}

func NewDispatcher(
	logger *slog.Logger,
	router *proxy.Router,
	sessions *session.Manager,
	dialogs *dialog.Manager,
	cfg config.Config,
	mediaEngine *media.Engine,
	cdrNotifier *cdr.Notifier,
) *Dispatcher {
	var digest *auth.Digest
	if cfg.AuthEnabled {
		digest = auth.NewDigest(cfg.AuthRealm, auth.Credentials(cfg.SIPUsers))
	}
	return &Dispatcher{
		logger:   logger,
		router:   router,
		sessions: sessions,
		dialogs:  dialogs,
		cfg:      cfg,
		digest:   digest,
		media:    mediaEngine,
		cdr:      cdrNotifier,
	}
}

func (d *Dispatcher) Handle(pkt Packet) error {
	if pkt.Msg.IsRequest {
		return d.handleRequest(pkt)
	}
	return d.handleResponse(pkt)
}

func (d *Dispatcher) handleRequest(pkt Packet) error {
	msg := pkt.Msg
	switch msg.Method {
	case sip.MethodOptions:
		return pkt.Reply(sip.BuildResponse(msg, 200, "OK", map[string]string{
			"Allow":     "INVITE, ACK, CANCEL, BYE, OPTIONS, REGISTER",
			"Accept":    "application/sdp",
			"Supported": "path, outbound",
		}))
	case sip.MethodRegister:
		return d.handleRegister(pkt)
	case sip.MethodInvite:
		return d.handleInvite(pkt)
	case sip.MethodAck:
		return d.handleAck(pkt)
	case sip.MethodBye:
		return d.handleBye(pkt)
	case sip.MethodCancel:
		return d.handleCancel(pkt)
	case sip.MethodNotify:
		return pkt.Reply(sip.BuildResponse(msg, 200, "OK", nil))
	default:
		return pkt.Reply(sip.BuildResponse(msg, 501, "Not Implemented", map[string]string{
			"Allow": "INVITE, ACK, CANCEL, BYE, OPTIONS, REGISTER",
		}))
	}
}

func (d *Dispatcher) handleRegister(pkt Packet) error {
	msg := pkt.Msg
	aor := sip.ExtractAOR(msg.Header("to"))
	if aor == "" {
		aor = sip.ExtractAOR(msg.Header("from"))
	}
	username := sip.ExtractUser(aor)

	if d.digest != nil {
		authHeader := msg.Header("authorization")
		user, ok := d.digest.Validate("REGISTER", authHeader)
		if !ok {
			d.logger.Info("REGISTER challenge", "aor", aor, "has_auth", authHeader != "")
			return pkt.Reply(sip.BuildResponse(msg, 401, "Unauthorized", map[string]string{
				"WWW-Authenticate": d.digest.ChallengeHeader(),
			}))
		}
		if username != "" && user != username {
			d.logger.Warn("REGISTER user mismatch", "uri_user", username, "auth_user", user)
			return pkt.Reply(sip.BuildResponse(msg, 403, "Forbidden", nil))
		}
		username = user
	}

	contact := msg.Header("contact")
	expires := sip.ContactExpires(contact, msg.Header("expires"), 3600)

	if contact == "" || expires == 0 || contact == "*" {
		d.router.Unregister(aor)
		d.logger.Info("REGISTER unregister", "aor", aor, "user", username)
		return pkt.Reply(sip.BuildResponse(msg, 200, "OK", map[string]string{"Expires": "0"}))
	}

	d.router.Register(proxy.Location{
		AOR:     aor,
		Contact: contact,
		Expires: expires,
	})
	d.logger.Info("REGISTER ok",
		"aor", aor,
		"user", username,
		"contact", contact,
		"expires", expires,
		"bindings", d.router.Count(),
	)
	return pkt.Reply(sip.BuildResponse(msg, 200, "OK", map[string]string{
		"Contact": fmt.Sprintf("%s;expires=%d", contact, expires),
		"Expires": fmt.Sprintf("%d", expires),
	}))
}

func (d *Dispatcher) handleInvite(pkt Packet) error {
	msg := pkt.Msg
	callID := msg.Header("call-id")
	target := msg.RequestURI
	if target == "" {
		target = msg.Header("to")
	}

	loc, ok := d.router.LookupFlexible(target)
	if !ok {
		loc, ok = d.router.LookupFlexible(sip.ExtractAOR(msg.Header("to")))
	}
	if !ok {
		d.logger.Info("INVITE 404", "target", target, "from", msg.Header("from"))
		return pkt.Reply(sip.BuildResponse(msg, 404, "Not Found", nil))
	}

	calleeAddr, err := sip.UDPAddrFromURI(loc.Contact)
	if err != nil {
		d.logger.Warn("contact inválido", "contact", loc.Contact, "err", err)
		return pkt.Reply(sip.BuildResponse(msg, 480, "Temporarily Unavailable", nil))
	}

	if err := pkt.Reply(sip.BuildResponse(msg, 100, "Trying", nil)); err != nil {
		return err
	}

	fwd := msg
	if d.cfg.MediaEnabled && d.media != nil && strings.TrimSpace(msg.Body) != "" {
		bridge, err := d.media.Open(callID)
		if err != nil {
			d.logger.Error("rtp bridge failed", "err", err)
			return pkt.Reply(sip.BuildResponse(msg, 500, "Server Internal Error", nil))
		}
		if audio, ok := media.ParseAudio(msg.Body); ok {
			bridge.SetCallerRemote(audio.IP, audio.Port)
		}
		fwd.Body = media.RewriteAudio(msg.Body, d.media.AdvertiseHost(), bridge.CalleeRTPPort())
	}

	branch := sip.NewBranch(callID + fmt.Sprintf("%d", time.Now().UnixNano()%100000))
	topVia := fmt.Sprintf("SIP/2.0/UDP %s;branch=%s;rport", d.cfg.AdvertisedAddr(), branch)
	requestURI := sip.ExtractURI(loc.Contact)
	forwarded := sip.ForwardRequest(fwd, requestURI, topVia)

	leg := &dialog.Leg{
		CallID:     callID,
		State:      dialog.StateEarly,
		CallerAddr: cloneUDPAddr(pkt.Remote),
		CalleeAddr: calleeAddr,
		CalleeURI:  requestURI,
		From:       msg.Header("from"),
		To:         msg.Header("to"),
		Branch:     branch,
	}
	d.dialogs.Put(leg)
	d.sessions.Upsert(&session.Call{
		ID:    callID,
		From:  sip.ExtractAOR(msg.Header("from")),
		To:    sip.ExtractAOR(msg.Header("to")),
		State: session.StateRinging,
	})

	d.logger.Info("INVITE proxy",
		"call_id", callID,
		"from", leg.From,
		"to", loc.AOR,
		"callee", calleeAddr.String(),
		"media", d.cfg.MediaEnabled,
	)

	return pkt.SendTo(forwarded, calleeAddr)
}

func (d *Dispatcher) handleAck(pkt Packet) error {
	msg := pkt.Msg
	leg, ok := d.dialogs.Get(msg.Header("call-id"))
	if !ok || leg.CalleeAddr == nil {
		d.logger.Debug("ACK sem diálogo", "call_id", msg.Header("call-id"))
		return nil
	}

	branch := sip.NewBranch(msg.Header("call-id") + "-ack")
	topVia := fmt.Sprintf("SIP/2.0/UDP %s;branch=%s;rport", d.cfg.AdvertisedAddr(), branch)
	requestURI := leg.CalleeURI
	if requestURI == "" {
		requestURI = msg.RequestURI
	}
	forwarded := sip.ForwardRequest(msg, requestURI, topVia)
	d.logger.Info("ACK proxy", "call_id", leg.CallID, "to", leg.CalleeAddr.String())
	return pkt.SendTo(forwarded, leg.CalleeAddr)
}

func (d *Dispatcher) handleBye(pkt Packet) error {
	msg := pkt.Msg
	leg, ok := d.dialogs.Get(msg.Header("call-id"))
	if !ok {
		return pkt.Reply(sip.BuildResponse(msg, 481, "Call/Transaction Does Not Exist", nil))
	}

	branch := sip.NewBranch(msg.Header("call-id") + "-bye")
	topVia := fmt.Sprintf("SIP/2.0/UDP %s;branch=%s;rport", d.cfg.AdvertisedAddr(), branch)

	var target *net.UDPAddr
	var requestURI string
	if sameEndpoint(pkt.Remote, leg.CallerAddr) {
		target = leg.CalleeAddr
		requestURI = leg.CalleeURI
	} else {
		target = leg.CallerAddr
		requestURI = msg.RequestURI
		if requestURI == "" {
			requestURI = "sip:" + sip.ExtractAOR(leg.From)
		}
	}

	if target == nil {
		return pkt.Reply(sip.BuildResponse(msg, 480, "Temporarily Unavailable", nil))
	}

	forwarded := sip.ForwardRequest(msg, requestURI, topVia)
	if err := pkt.SendTo(forwarded, target); err != nil {
		return err
	}

	// end locally; final 200 from peer is still relayed via handleResponse
	d.finishCall(leg.CallID)
	leg.State = dialog.StateTerminated
	d.logger.Info("BYE proxy", "call_id", leg.CallID, "to", target.String())
	return nil
}

func (d *Dispatcher) handleCancel(pkt Packet) error {
	msg := pkt.Msg
	leg, ok := d.dialogs.Get(msg.Header("call-id"))
	if !ok || leg.CalleeAddr == nil {
		return pkt.Reply(sip.BuildResponse(msg, 481, "Call/Transaction Does Not Exist", nil))
	}
	_ = pkt.Reply(sip.BuildResponse(msg, 200, "OK", nil))

	branch := sip.NewBranch(msg.Header("call-id") + "-cancel")
	topVia := fmt.Sprintf("SIP/2.0/UDP %s;branch=%s;rport", d.cfg.AdvertisedAddr(), branch)
	forwarded := sip.ForwardRequest(msg, leg.CalleeURI, topVia)
	if d.media != nil {
		d.media.Close(leg.CallID)
	}
	d.finishCall(leg.CallID)
	d.logger.Info("CANCEL proxy", "call_id", leg.CallID)
	return pkt.SendTo(forwarded, leg.CalleeAddr)
}

func (d *Dispatcher) handleResponse(pkt Packet) error {
	msg := pkt.Msg
	callID := msg.Header("call-id")
	leg, ok := d.dialogs.Get(callID)
	if !ok {
		d.logger.Debug("resposta SIP sem diálogo", "call_id", callID, "status", msg.StatusCode)
		return nil
	}

	// Responses from callee go to caller (and vice-versa for BYE responses).
	var target *net.UDPAddr
	if sameEndpoint(pkt.Remote, leg.CalleeAddr) {
		target = leg.CallerAddr
	} else if sameEndpoint(pkt.Remote, leg.CallerAddr) {
		target = leg.CalleeAddr
	} else {
		// fallback: treat as callee response
		target = leg.CallerAddr
	}
	if target == nil {
		return nil
	}

	fwd := msg
	switch {
	case msg.StatusCode >= 180 && msg.StatusCode < 200:
		leg.State = dialog.StateEarly
		if call, ok := d.sessions.Get(callID); ok {
			call.State = session.StateRinging
			d.sessions.Upsert(call)
		}
		fwd = d.rewriteAnswerSDP(callID, msg)
	case msg.StatusCode >= 200 && msg.StatusCode < 300:
		cseq := strings.ToUpper(msg.Header("cseq"))
		if strings.Contains(cseq, "INVITE") {
			leg.State = dialog.StateConfirmed
			d.sessions.MarkAnswered(callID)
			if contact := msg.Header("contact"); contact != "" {
				leg.CalleeURI = sip.ExtractURI(contact)
				if addr, err := sip.UDPAddrFromURI(contact); err == nil {
					leg.CalleeAddr = addr
				}
			}
			fwd = d.rewriteAnswerSDP(callID, msg)
		}
		if strings.Contains(cseq, "BYE") {
			leg.State = dialog.StateTerminated
			d.finishCall(callID)
			d.dialogs.Delete(callID)
		}
	case msg.StatusCode >= 300:
		d.finishCall(callID)
		leg.State = dialog.StateTerminated
	}

	forwarded := sip.ForwardResponse(fwd)
	d.logger.Info("SIP response proxy",
		"call_id", callID,
		"status", msg.StatusCode,
		"to", target.String(),
	)
	return pkt.SendTo(forwarded, target)
}

func (d *Dispatcher) rewriteAnswerSDP(callID string, msg sip.Message) sip.Message {
	if !d.cfg.MediaEnabled || d.media == nil || strings.TrimSpace(msg.Body) == "" {
		return msg
	}
	bridge, ok := d.media.Get(callID)
	if !ok {
		return msg
	}
	if audio, ok := media.ParseAudio(msg.Body); ok {
		bridge.SetCalleeRemote(audio.IP, audio.Port)
	}
	out := msg
	out.Body = media.RewriteAudio(msg.Body, d.media.AdvertiseHost(), bridge.CallerRTPPort())
	return out
}

func (d *Dispatcher) finishCall(callID string) {
	call, newlyEnded := d.sessions.MarkEnded(callID)
	if d.media != nil {
		d.media.Close(callID)
	}
	if newlyEnded && call != nil && d.cdr != nil {
		d.cdr.NotifyEnded(*call)
	}
}

func sameEndpoint(a, b *net.UDPAddr) bool {
	if a == nil || b == nil {
		return false
	}
	return a.IP.Equal(b.IP) && a.Port == b.Port
}

func cloneUDPAddr(addr *net.UDPAddr) *net.UDPAddr {
	if addr == nil {
		return nil
	}
	ip := make(net.IP, len(addr.IP))
	copy(ip, addr.IP)
	return &net.UDPAddr{IP: ip, Port: addr.Port, Zone: addr.Zone}
}

func (d *Dispatcher) Digest() *auth.Digest {
	return d.digest
}