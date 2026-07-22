package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/kvoip/kvoip-pbx/internal/config"
	"github.com/kvoip/kvoip-pbx/internal/handlers"
	"github.com/kvoip/kvoip-pbx/internal/sip"
)

// Server is the SIP network listener facade.
type Server struct {
	cfg        config.Config
	logger     *slog.Logger
	dispatcher *handlers.Dispatcher
}

func New(cfg config.Config, logger *slog.Logger, dispatcher *handlers.Dispatcher) *Server {
	return &Server{
		cfg:        cfg,
		logger:     logger,
		dispatcher: dispatcher,
	}
}

// ListenAndServeUDP starts the UDP SIP listener until context cancellation.
func (s *Server) ListenAndServeUDP(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", s.cfg.ListenAddr())
	if err != nil {
		return fmt.Errorf("resolver endereço UDP: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("listen UDP: %w", err)
	}
	defer conn.Close()

	s.logger.Info("listener SIP UDP iniciado", "addr", addr.String())

	buffer := make([]byte, s.cfg.BufferSize)
	errCh := make(chan error, 1)

	go func() {
		for {
			n, remote, readErr := conn.ReadFromUDP(buffer)
			if readErr != nil {
				select {
				case <-ctx.Done():
					return
				default:
					errCh <- readErr
					return
				}
			}

			payload := make([]byte, n)
			copy(payload, buffer[:n])
			remoteCopy := cloneAddr(remote)

			msg := sip.Parse(payload)
			s.logger.Debug("datagrama SIP recebido",
				"remote", remoteCopy.String(),
				"bytes", n,
				"method", string(msg.Method),
				"status", msg.StatusCode,
				"request", msg.IsRequest,
			)

			pkt := handlers.Packet{
				Msg:    msg,
				Remote: remoteCopy,
				Reply: func(data []byte) error {
					_, err := conn.WriteToUDP(data, remoteCopy)
					return err
				},
				SendTo: func(data []byte, addr *net.UDPAddr) error {
					_, err := conn.WriteToUDP(data, addr)
					return err
				},
			}

			if handleErr := s.dispatcher.Handle(pkt); handleErr != nil {
				s.logger.Warn("falha ao processar SIP", "err", handleErr)
			}
		}
	}()

	select {
	case <-ctx.Done():
		_ = conn.Close()
		s.logger.Info("encerrando listener SIP")
		return nil
	case err := <-errCh:
		return fmt.Errorf("leitura UDP: %w", err)
	}
}

func cloneAddr(addr *net.UDPAddr) *net.UDPAddr {
	if addr == nil {
		return nil
	}
	ip := make(net.IP, len(addr.IP))
	copy(ip, addr.IP)
	return &net.UDPAddr{IP: ip, Port: addr.Port, Zone: addr.Zone}
}
