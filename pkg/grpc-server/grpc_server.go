package grpcserver

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	gen "github.com/Ranik23/avito-tech-spring/api/proto/gen/pvz_v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
	logger *slog.Logger
	config *Config
}


type Config struct {
	Host 				string
	Port 				string	
	StartMsg			string
	ShutdownTimeout   	time.Duration
}

// убрать srv gen.PVZServiceServer и поставить что-то более общее
func New(logger *slog.Logger,  cfg *Config, srv gen.PVZServiceServer) *Server {
	grpcServer := grpc.NewServer()

	gen.RegisterPVZServiceServer(grpcServer, srv)

	return &Server{
		server: grpcServer,
		logger: logger,
		config: cfg,
	}
}


func (s *Server) Start(ctx context.Context) error {
	s.logger.Info(s.config.StartMsg)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func () error {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		select{
		case <-ctx.Done():
			s.logger.Info("Context Cancelled")
		case <-sigCh:
			s.logger.Info("Signal Detected!")
		}

		s.logger.Info("Gracefully shutting down GRPC Server")
		s.server.GracefulStop()

		return nil
	})


	g.Go(func() error {
		listener, err := net.Listen("tcp", s.config.Host + ":" + s.config.Port)
		if err != nil {
			return err
		}

		if err := s.server.Serve(listener); err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
			} else {
				return err
			}
		}
		return nil
	})


	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

