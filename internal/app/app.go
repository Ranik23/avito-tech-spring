package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/Ranik23/avito-tech-spring/internal/config"
	grpccontrollers "github.com/Ranik23/avito-tech-spring/internal/controllers/grpc"
	httpcontrollers "github.com/Ranik23/avito-tech-spring/internal/controllers/http"
	"github.com/Ranik23/avito-tech-spring/internal/hasher"
	"github.com/Ranik23/avito-tech-spring/internal/repository/postgresql"
	"github.com/Ranik23/avito-tech-spring/internal/service"
	"github.com/Ranik23/avito-tech-spring/internal/token"
	grpcserver "github.com/Ranik23/avito-tech-spring/pkg/grpc-server"
	httpserver "github.com/Ranik23/avito-tech-spring/pkg/http-server"
	"github.com/gin-gonic/gin"
	"github.com/lmittmann/tint"
	"golang.org/x/sync/errgroup"
)

type App struct {
	service service.Service
	logger  *slog.Logger

	cfg *config.Config

	httpServer *httpserver.Server
	grcpServer *grpcserver.Server
}

func NewApp() (*App, error) {
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	
	logger.Info("Loading configuration...")

	cfg, err := config.LoadConfig("config/", ".env")
	if err != nil {
		logger.Error("Failed to load config", slog.String("error", err.Error()))
		return nil, err
	}
	logger.Info("Configuration loaded")

	logger.Info("Connecting to database...")
	pool, err := cfg.Storage.Connect()
	if err != nil {
		logger.Error("Failed to connect to database", slog.String("error", err.Error()))
		return nil, err
	}
	logger.Info("Connected to database")

	ctxManager := postgresql.NewCtxManager(pool)
	txManager := postgresql.NewTxManager(pool, logger, ctxManager)

	logger.Info("Initializing repositories...")
	userRepo := postgresql.NewPostgresUserRepository(ctxManager, logger)
	productRepo := postgresql.NewPostgresProductRepository(ctxManager, logger)
	pvzRepo := postgresql.NewPostgresPvzRepository(ctxManager, logger)
	receptionRepo := postgresql.NewPostgresReceptionRepository(ctxManager, logger)

	logger.Info("Initializing services...")
	tokenService := token.NewToken(cfg.SecretKey)
	hasher := hasher.NewHasher()

	authService := service.NewAuthService(userRepo, txManager, tokenService, hasher, logger)
	pvzService := service.NewPVZService(pvzRepo, receptionRepo, cfg.Cities, productRepo, txManager, logger)
	service := service.NewService(authService, pvzService)

	logger.Info("Initializing controllers...")
	authController := httpcontrollers.NewAuthController(service, logger)
	pvzController := httpcontrollers.NewPVZController(service)

	httpServerConfig := &httpserver.Config{
		Port: cfg.HTTPServer.Port,
		Host: cfg.HTTPServer.Host,
		StartMsg: "Hello, I Am A HTTP Server",
	}

	grpcServerConfig := &grpcserver.Config{
		Host: cfg.GRPCServer.Host,
		Port: cfg.GRPCServer.Port,
		StartMsg: "Hello, I Am A GRPC Server",
	}

	grpcServerImpl := grpccontrollers.NewPVZServer(service)


	logger.Info("Setting up HTTP routes...")
	router := gin.New()
	SetUpRoutes(router, authController, pvzController, tokenService)

	logger.Info("Creating HTTP server...")
	httpServer := httpserver.New(logger, httpServerConfig, router)

	logger.Info("Creating GRPC server...")
	grpcServer := grpcserver.New(logger, grpcServerConfig, grpcServerImpl)


	logger.Info("App initialization complete")

	return &App{
		service:    service,
		logger:     logger,
		cfg:        cfg,
		httpServer: httpServer,
		grcpServer: grpcServer,
	}, nil
}

func (a *App) Start() error {
	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		if err := a.httpServer.Start(context.TODO()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("Failed to start HTTP server", slog.String("error", err.Error()))
			return err
		}
		return nil 
	})

	g.Go(func() error {
		if err := a.grcpServer.Start(context.TODO()); err != nil {
			a.logger.Error("Failed to start GRPC Server", slog.String("error", err.Error()))
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	a.logger.Info("HTTP Server Gracefully stopped")
	a.logger.Info("GRPC Server Gracefully stopped")

	
	return nil
}
