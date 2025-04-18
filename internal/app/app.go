package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Ranik23/avito-tech-spring/internal/config"
	httpcontrollers "github.com/Ranik23/avito-tech-spring/internal/controllers/http"
	"github.com/Ranik23/avito-tech-spring/internal/hasher"
	"github.com/Ranik23/avito-tech-spring/internal/repository/postgresql"
	"github.com/Ranik23/avito-tech-spring/internal/service"
	"github.com/Ranik23/avito-tech-spring/internal/token"
	httpserver "github.com/Ranik23/avito-tech-spring/pkg/http-server"
	"github.com/gin-gonic/gin"
	grpcserver "google.golang.org/grpc"
)

type App struct {
	service service.Service
	logger  *slog.Logger

	cfg *config.Config

	httpServer *httpserver.Server
	grcpServer *grpcserver.Server
}

func NewApp() (*App, error) {
	logger := slog.Default()
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
	userRepo := postgresql.NewPostgresUserRepository(ctxManager)
	productRepo := postgresql.NewPostgresProductRepository(ctxManager)
	pvzRepo := postgresql.NewPostgresPvzRepository(ctxManager)
	receptionRepo := postgresql.NewPostgresReceptionRepository(ctxManager)

	logger.Info("Initializing services...")
	tokenService := token.NewToken(cfg.SecretKey)
	hasher := hasher.NewHasher()

	authService := service.NewAuthService(userRepo, txManager, tokenService, hasher, logger)
	pvzService := service.NewPVZService(pvzRepo, receptionRepo, cfg.Cities, productRepo, txManager, logger)
	service := service.NewService(authService, pvzService)

	logger.Info("Initializing controllers...")
	authController := httpcontrollers.NewAuthController(service)
	pvzController := httpcontrollers.NewPVZController(service)

	httpServerConfig := httpserver.Config{
		Port: cfg.HTTPServer.Port,
		Host: cfg.HTTPServer.Host,
	}

	logger.Info("Setting up HTTP routes...")
	router := gin.New()
	SetUpRoutes(router, authController, pvzController, tokenService)

	logger.Info("Creating HTTP server...")
	httpServer := httpserver.New(logger, httpServerConfig, router)

	logger.Info("App initialization complete")

	return &App{
		service:    service,
		logger:     logger,
		cfg:        cfg,
		httpServer: httpServer,
	}, nil
}

func (a *App) Start() error {
	a.logger.Info("Starting the HTTP server...",
		slog.String("host", a.cfg.HTTPServer.Host),
		slog.String("port", a.cfg.HTTPServer.Port),
	)

	if err := a.httpServer.Start(context.TODO()); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.logger.Error("Failed to start HTTP server", slog.String("error", err.Error()))
		return err
	}

	a.logger.Info("HTTP server gracefully stopped")
	return nil
}
