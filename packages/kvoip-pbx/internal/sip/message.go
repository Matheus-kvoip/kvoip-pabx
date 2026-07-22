package sip

import (
	"fmt"
	"strconv"
	"strings"
)

// Method is a SIP request method.
type Method string

const (
	MethodInvite   Method = "INVITE"
	MethodAck      Method = "ACK"
	MethodBye      Method = "BYE"
	MethodCancel   Method = "CANCEL"
	MethodOptions  Method = "OPTIONS"
	MethodRegister Method = "REGISTER"
	MethodNotify   Method = "NOTIFY"
)

// Message is a minimal SIP message representation.
type Message struct {
	Raw        string
	StartLine  string
	Method     Method
	RequestURI string
	IsRequest  bool
	StatusCode int
	Reason     string
	Vias       []string
	Headers    map[string]string
	Body       string
}

// Parse extracts a lightweight view of a SIP datagram.
func Parse(raw []byte) Message {
	text := string(raw)
	msg := Message{
		Raw:     text,
		Headers: map[string]string{},
	}

	parts := strings.SplitN(text, "\r\n\r\n", 2)
	headerBlock := parts[0]
	if len(parts) == 2 {
		msg.Body = parts[1]
	}

	lines := strings.Split(headerBlock, "\r\n")
	if len(lines) == 0 {
		return msg
	}

	msg.StartLine = strings.TrimSpace(lines[0])
	fields := strings.Fields(msg.StartLine)
	if len(fields) >= 1 && strings.HasPrefix(fields[0], "SIP/") {
		msg.IsRequest = false
		if len(fields) >= 2 {
			code, _ := strconv.Atoi(fields[1])
			msg.StatusCode = code
		}
		if len(fields) >= 3 {
			msg.Reason = strings.Join(fields[2:], " ")
		}
	} else if len(fields) >= 1 {
		msg.IsRequest = true
		msg.Method = Method(strings.ToUpper(fields[0]))
		if len(fields) >= 2 {
			msg.RequestURI = fields[1]
		}
	}

	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		name, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(name))
		val := strings.TrimSpace(value)
		if key == "via" {
			msg.Vias = append(msg.Vias, val)
			continue
		}
		if _, exists := msg.Headers[key]; !exists {
			msg.Headers[key] = val
		}
	}

	return msg
}

func (m Message) Header(name string) string {
	return m.Headers[strings.ToLower(name)]
}

// BuildResponse creates a SIP response for the given request.
func BuildResponse(req Message, status int, reason string, extra map[string]string) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "SIP/2.0 %d %s\r\n", status, reason)

	for _, via := range req.Vias {
		fmt.Fprintf(&b, "Via: %s\r\n", via)
	}
	writeCopiedDialogHeaders(&b, req)

	for name, value := range extra {
		fmt.Fprintf(&b, "%s: %s\r\n", name, value)
	}

	fmt.Fprintf(&b, "User-Agent: KVOIP-PBX/0.1\r\n")
	fmt.Fprintf(&b, "Content-Length: 0\r\n\r\n")
	return []byte(b.String())
}

// ForwardRequest builds a proxied request with a new top Via and Request-URI.
func ForwardRequest(req Message, requestURI, topVia string) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s SIP/2.0\r\n", req.Method, requestURI)
	fmt.Fprintf(&b, "Via: %s\r\n", topVia)
	for _, via := range req.Vias {
		fmt.Fprintf(&b, "Via: %s\r\n", via)
	}

	ordered := []string{"from", "to", "call-id", "cseq", "contact", "max-forwards", "content-type"}
	seen := map[string]bool{"via": true}
	for _, key := range ordered {
		if value := req.Header(key); value != "" {
			if key == "max-forwards" {
				if n, err := strconv.Atoi(value); err == nil && n > 0 {
					value = strconv.Itoa(n - 1)
				}
			}
			fmt.Fprintf(&b, "%s: %s\r\n", headerDisplayName(key), value)
			seen[key] = true
		}
	}
	for key, value := range req.Headers {
		if seen[key] {
			continue
		}
		if key == "content-length" {
			continue
		}
		fmt.Fprintf(&b, "%s: %s\r\n", headerDisplayName(key), value)
	}

	fmt.Fprintf(&b, "Content-Length: %d\r\n\r\n", len(req.Body))
	b.WriteString(req.Body)
	return []byte(b.String())
}

// ForwardResponse strips the top Via and relays the response upstream.
func ForwardResponse(res Message) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "SIP/2.0 %d %s\r\n", res.StatusCode, res.Reason)

	vias := res.Vias
	if len(vias) > 0 {
		vias = vias[1:]
	}
	for _, via := range vias {
		fmt.Fprintf(&b, "Via: %s\r\n", via)
	}

	ordered := []string{"from", "to", "call-id", "cseq", "contact", "content-type"}
	seen := map[string]bool{"via": true}
	for _, key := range ordered {
		if value := res.Header(key); value != "" {
			fmt.Fprintf(&b, "%s: %s\r\n", headerDisplayName(key), value)
			seen[key] = true
		}
	}
	for key, value := range res.Headers {
		if seen[key] || key == "content-length" {
			continue
		}
		fmt.Fprintf(&b, "%s: %s\r\n", headerDisplayName(key), value)
	}

	fmt.Fprintf(&b, "Content-Length: %d\r\n\r\n", len(res.Body))
	b.WriteString(res.Body)
	return []byte(b.String())
}

func writeCopiedDialogHeaders(b *strings.Builder, req Message) {
	if from := req.Header("from"); from != "" {
		fmt.Fprintf(b, "From: %s\r\n", from)
	}
	to := req.Header("to")
	if to != "" {
		if !strings.Contains(strings.ToLower(to), ";tag=") {
			to = to + ";tag=kvoip-" + shortTag(req.Header("call-id"))
		}
		fmt.Fprintf(b, "To: %s\r\n", to)
	}
	if callID := req.Header("call-id"); callID != "" {
		fmt.Fprintf(b, "Call-ID: %s\r\n", callID)
	}
	if cseq := req.Header("cseq"); cseq != "" {
		fmt.Fprintf(b, "CSeq: %s\r\n", cseq)
	}
}

func headerDisplayName(key string) string {
	switch key {
	case "call-id":
		return "Call-ID"
	case "cseq":
		return "CSeq"
	case "max-forwards":
		return "Max-Forwards"
	case "content-type":
		return "Content-Type"
	case "user-agent":
		return "User-Agent"
	default:
		if key == "" {
			return key
		}
		return strings.ToUpper(key[:1]) + key[1:]
	}
}

func shortTag(seed string) string {
	if seed == "" {
		return "pbx"
	}
	seed = strings.ReplaceAll(seed, "@", "")
	seed = strings.ReplaceAll(seed, ".", "")
	if len(seed) > 8 {
		seed = seed[:8]
	}
	return seed
}
