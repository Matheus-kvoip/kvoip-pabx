package transport

// Kind identifies a SIP transport.
type Kind string

const (
	UDP Kind = "udp"
	TCP Kind = "tcp"
	TLS Kind = "tls"
)

// Endpoint is a remote SIP peer address.
type Endpoint struct {
	Host string
	Port int
	Kind Kind
}
