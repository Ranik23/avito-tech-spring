package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	gen "github.com/Ranik23/avito-tech-spring/api/proto/gen/pvz_v1"
	"github.com/Ranik23/avito-tech-spring/internal/config"
	grpccontrollers "github.com/Ranik23/avito-tech-spring/internal/controllers/grpc"
	httpcontroll "github.com/Ranik23/avito-tech-spring/internal/controllers/http"
	"github.com/Ranik23/avito-tech-spring/internal/hasher"
	"github.com/Ranik23/avito-tech-spring/internal/repository/postgresql"
	"github.com/Ranik23/avito-tech-spring/internal/service"
	"github.com/Ranik23/avito-tech-spring/internal/token"
	"github.com/Ranik23/avito-tech-spring/pkg/closure"
	grpcserver "github.com/Ranik23/avito-tech-spring/pkg/grpc-server"
	httpserver "github.com/Ranik23/avito-tech-spring/pkg/http-server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lmittmann/tint"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	service 		service.Service
	logger  		*slog.Logger

	cfg 			*config.Config

	httpServer 		*httpserver.Server
	grcpServer 		*grpcserver.Server
	gatewayServer 	*httpserver.Server

	closer 			*closure.Closer
}

func NewApp() (*App, error) {
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	
	logger.Info("Loading configuration...")

	closer := closure.NewCloser()

	logger.Info("Closure initialized!")

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
	closer.Add(func(ctx context.Context) error {
		logger.Info("Closing Pool!")
		pool.Close()
		return nil
	})

	config := NewCORSConfig()

	logger.Info("Connected to database")

	ctxManager := postgresql.NewCtxManager(pool)
	txManager := postgresql.NewTxManager(pool, logger, ctxManager)

	logger.Info("Initializing repositories...")
	userRepo := postgresql.NewPostgresUserRepository(ctxManager, logger)
	productRepo := postgresql.NewPostgresProductRepository(ctxManager, logger)
	pvzRepo := postgresql.NewPostgresPvzRepository(ctxManager, logger)
	receptionRepo := postgresql.NewPostgresReceptionRepository(ctxManager, logger)

	logger.Info("Initializing services...")
	tokenService := token.NewToken(cfg.SecretKey, logger)
	passwordhasher := hasher.NewHasher()

	authService := service.NewAuthService(userRepo, txManager, tokenService, passwordhasher, logger)
	pvzService := service.NewPVZService(pvzRepo, receptionRepo, cfg.Cities, productRepo, txManager, logger)
	service := service.NewService(authService, pvzService)

	logger.Info("Initializing controllers...")
	authController := httpcontroll.NewAuthController(service, logger)
	pvzController := httpcontroll.NewPVZController(service, logger)





	logger.Info("Creating Gateway Server")

	gateWayConfig := &httpserver.Config{
		Port: 		cfg.GatewayServer.Port,
		Host: 		cfg.GatewayServer.Host,
		StartMsg: 	"Hello, I am A Gateway Server",
	}

	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	grpcAddr := fmt.Sprintf("%s:%s", cfg.GRPCServer.Host, cfg.GRPCServer.Port)

	err = gen.RegisterPVZServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		log.Printf("Failed to register gateway: %v", err)
		return nil, err
	}

	gatewayServer := httpserver.New(logger, gateWayConfig, mux)

	logger.Info("Gateway Server Created")





	logger.Info("Creating HTTP server...")

	httpServerConfig := &httpserver.Config{
		Port: 		cfg.HTTPServer.Port,
		Host: 		cfg.HTTPServer.Host,
		StartMsg: 	"Hello, I Am A HTTP Server",
	}

	logger.Info("Setting up HTTP routes...")
	router := gin.New()

	router.Use(cors.New(config))

	SetUpRoutes(router, authController, pvzController, tokenService)

	httpServer := httpserver.New(logger, httpServerConfig, router)

	logger.Info("HTTP Server Created")






	logger.Info("Creating GRPC server...")

	grpcServerConfig := &grpcserver.Config{
		Host: 		cfg.GRPCServer.Host,
		Port: 		cfg.GRPCServer.Port,
		StartMsg: 	"Hello, I Am A GRPC Server",
	}

	grpcServerImpl := grpccontrollers.NewPVZServer(service)
	
	grpcServer := grpc.NewServer()
	gen.RegisterPVZServiceServer(grpcServer, grpcServerImpl)

	grpcSrv := grpcserver.New(logger, grpcServerConfig, grpcServer)

	logger.Info("GRCP Server Created")






	logger.Info("App initialization complete")
	return &App{
		service:    service,
		logger:     logger,
		cfg:        cfg,
		httpServer: httpServer,
		grcpServer: grpcSrv,
		closer: 	closer,
		gatewayServer: gatewayServer,
	}, nil
}

func (a *App) Start(ctx context.Context) error {

	defer func() {
		if err := a.closer.Close(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		if err := a.gatewayServer.Start(ctx); err != nil && errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("Failed to start Gateway Server", slog.String("error", err.Error()))
			return err
		}
		return nil
	})

	g.Go(func() error {
		if err := a.httpServer.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("Failed to start HTTP server", slog.String("error", err.Error()))
			return err
		}
		return nil 
	})

	g.Go(func() error {
		if err := a.grcpServer.Start(ctx); err != nil {
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
	a.logger.Info("Gateway Server Gracefully stopped")

	
	return nil
}
