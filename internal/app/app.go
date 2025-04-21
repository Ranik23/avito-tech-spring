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
	"github.com/Ranik23/avito-tech-spring/internal/controllers/grpc/interceptors"
	httpcontrollers "github.com/Ranik23/avito-tech-spring/internal/controllers/http"
	"github.com/Ranik23/avito-tech-spring/internal/controllers/http/middleware"
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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	service 		service.Service
	logger  		*slog.Logger

	cfg 			*config.Config

	httpServer 		*httpserver.Server
	grcpServer 		*grpcserver.Server
	gatewayServer 	*httpserver.Server
	metricServer	*httpserver.Server

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
	authController := httpcontrollers.NewAuthController(service, logger)
	pvzController := httpcontrollers.NewPVZController(service, logger)


	gatewayServer, err := createGateWayServer(logger, cfg)
	if err != nil {
		logger.Error("Failed to create GateWay Server", slog.String("error", err.Error()))
		return nil, err
	}

	httpServer := createHTTPServer(logger, cfg, authController, pvzController, tokenService)
	grpcServer := createGRPCServer(logger, service, cfg)
	metricServer := createMetricsServer(logger, cfg)

	logger.Info("App initialization complete")

	return &App{
		service:    	service,
		logger:     	logger,
		cfg:        	cfg,
		httpServer: 	httpServer,
		grcpServer: 	grpcServer,
		closer: 		closer,
		gatewayServer: 	gatewayServer,
		metricServer: 	metricServer,
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

	g.Go(func() error {
		if err := a.metricServer.Start(ctx); err != nil {
			a.logger.Error("Failed to start Metric Server", slog.String("error", err.Error()))
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
	a.logger.Info("HTTP Metric Server Gracefully stopped")

	
	return nil
}


func createGRPCServer(logger *slog.Logger, service service.Service, cfg *config.Config) *grpcserver.Server {
	logger.Info("Creating GRPC server...")

	grpcServerConfig := &grpcserver.Config{
		Host: 		cfg.GRPCServer.Host,
		Port: 		cfg.GRPCServer.Port,
		StartMsg: 	"Hello, I Am A GRPC Server",
	}

	grpcServerImpl := grpccontrollers.NewPVZServer(service)
	
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
			interceptors.LoggingUnaryInterceptor(logger),
		),
	)

	gen.RegisterPVZServiceServer(grpcServer, grpcServerImpl)
	reflection.Register(grpcServer)
	
	grpcSrv := grpcserver.New(logger, grpcServerConfig, grpcServer)

	logger.Info("GRCP Server Created")

	return grpcSrv
}


func createGateWayServer(logger *slog.Logger, cfg *config.Config) (*httpserver.Server, error) {
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

	err := gen.RegisterPVZServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		log.Printf("Failed to register gateway: %v", err)
		return nil, err
	}
	
	gatewayServer := httpserver.New(logger, gateWayConfig, mux)

	logger.Info("Gateway Server Created")

	return gatewayServer, nil
}


func createHTTPServer(logger *slog.Logger, cfg *config.Config, authController httpcontrollers.AuthController, 
	pvzController httpcontrollers.PvzController, tokenService token.Token) *httpserver.Server {

	logger.Info("Creating HTTP server...")

	config := NewCORSConfig()

	httpServerConfig := &httpserver.Config{
		Port: 		cfg.HTTPServer.Port,
		Host: 		cfg.HTTPServer.Host,
		StartMsg: 	"Hello, I Am A HTTP Server",
	}

	logger.Info("Setting up HTTP routes...")

	router := gin.New()

	router.Use(middleware.Duration())
	router.Use(cors.New(config))

	SetUpRoutes(router, authController, pvzController, tokenService)

	httpServer := httpserver.New(logger, httpServerConfig, router)

	logger.Info("HTTP Server Created")

	return httpServer
}


func createMetricsServer(logger *slog.Logger, cfg *config.Config) *httpserver.Server {
	logger.Info("Creating HTTP Metric Server...")

	httpServerConfig := &httpserver.Config{
		Port: 		cfg.MetricServer.Port,
		Host: 		cfg.MetricServer.Host,
		StartMsg: 	"Hello, I Am A HTTP Metrics Server",
	}

	router := gin.New()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	httpServer := httpserver.New(logger, httpServerConfig, router)

	logger.Info("HTTP Metric Server Created")

	return httpServer
}