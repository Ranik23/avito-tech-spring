package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	logger *slog.Logger
	server *http.Server
	config *Config
}

type Config struct {
	Host			  string
	Port              string
	StartMsg          string
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func New(logger *slog.Logger, config *Config, handler http.Handler) *Server {
	server := &http.Server{
		Handler:           handler,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		Addr:              config.Host + ":" + config.Port,
	}

	s := Server{
		logger: logger,
		server: server,
		config: config,
	}

	return &s
}

func (a *Server) Start(ctx context.Context) error {
	a.logger.Info(a.config.StartMsg)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		select{
		case <-ctx.Done():
			a.logger.Info("Context Cancelled")
		case <- sigCh:
			a.logger.Info("Signal Detected")
		}
	
		ctx, cancel := context.WithTimeout(context.Background(), a.config.ShutdownTimeout)
		defer cancel()

		a.logger.Info("Gracefully Shutting Down HTTP Server")
		err := a.server.Shutdown(ctx)
		if err != nil {
			return err
		}

		return nil
	})


	g.Go(func() error {
		err := a.server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
			} else {
				return err
			}
		}
		return nil
	})

	err := g.Wait()
	if err != nil {
		return err
	}

	return nil
}