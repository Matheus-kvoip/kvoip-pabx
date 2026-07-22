package media

import (
	"fmt"
	"strconv"
	"strings"
)

// Audio describes the first audio media line in an SDP body.
type Audio struct {
	IP   string
	Port int
}

// ParseAudio extracts connection IP and audio port from SDP.
func ParseAudio(sdp string) (Audio, bool) {
	sdp = strings.ReplaceAll(sdp, "\r\n", "\n")
	var sessionIP string
	var mediaIP string
	var port int
	inAudio := false

	for _, line := range strings.Split(sdp, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "c=IN IP4 "):
			ip := strings.TrimSpace(strings.TrimPrefix(line, "c=IN IP4 "))
			if inAudio {
				mediaIP = ip
			} else {
				sessionIP = ip
			}
		case strings.HasPrefix(line, "m=audio "):
			inAudio = true
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				p, err := strconv.Atoi(fields[1])
				if err == nil {
					port = p
				}
			}
		case strings.HasPrefix(line, "m=") && !strings.HasPrefix(line, "m=audio "):
			inAudio = false
		}
	}

	ip := mediaIP
	if ip == "" {
		ip = sessionIP
	}
	if ip == "" || port <= 0 {
		return Audio{}, false
	}
	return Audio{IP: ip, Port: port}, true
}

// RewriteAudio forces connection IP and audio RTP port.
func RewriteAudio(sdp, host string, port int) string {
	if strings.TrimSpace(sdp) == "" {
		return sdp
	}
	lines := strings.Split(strings.ReplaceAll(sdp, "\r\n", "\n"), "\n")
	out := make([]string, 0, len(lines)+1)
	hasC := false

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if line == "" && len(out) > 0 {
			continue
		}
		switch {
		case strings.HasPrefix(line, "c=IN IP4 "):
			hasC = true
			out = append(out, fmt.Sprintf("c=IN IP4 %s", host))
		case strings.HasPrefix(line, "m=audio "):
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				fields[1] = strconv.Itoa(port)
				out = append(out, strings.Join(fields, " "))
			} else {
				out = append(out, line)
			}
		case strings.HasPrefix(line, "a=rtcp:"):
			continue
		default:
			out = append(out, line)
		}
	}

	if !hasC {
		inserted := false
		withC := make([]string, 0, len(out)+1)
		for _, line := range out {
			withC = append(withC, line)
			if !inserted && strings.HasPrefix(line, "o=") {
				withC = append(withC, fmt.Sprintf("c=IN IP4 %s", host))
				inserted = true
			}
		}
		if !inserted && len(withC) > 0 {
			withC = append([]string{withC[0], fmt.Sprintf("c=IN IP4 %s", host)}, withC[1:]...)
		}
		out = withC
	}

	body := strings.Join(out, "\r\n")
	if !strings.HasSuffix(body, "\r\n") {
		body += "\r\n"
	}
	return body
}
