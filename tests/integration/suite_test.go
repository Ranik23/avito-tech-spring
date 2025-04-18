//go:build integration

package integration

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ranik23/avito-tech-spring/internal/config"
	"github.com/Ranik23/avito-tech-spring/internal/hasher"
	"github.com/Ranik23/avito-tech-spring/internal/repository"
	"github.com/Ranik23/avito-tech-spring/internal/repository/manager"
	"github.com/Ranik23/avito-tech-spring/internal/service"
	"github.com/Ranik23/avito-tech-spring/internal/token"
	"github.com/Ranik23/avito-tech-spring/tests/integration/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type TestSuite struct {
	suite.Suite
	psqlContainer  *util.PostgreSQLContainer
	server         *httptest.Server
	grpcClient     grpc.ClientConnInterface
	ctxManager     manager.CtxManager
	pool           *pgxpool.Pool

	authService service.AuthService
	pvzService  service.PVZService
}

func (s *TestSuite) SetupSuite() {
	cfg, err := config.LoadConfig("../../config/", "../../.env")
	s.Require().NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	psqlContainer, err := util.NewPostgreSQLContainer(ctx)
	s.Require().NoError(err)

	s.psqlContainer = psqlContainer

	err = util.RunMigrations(psqlContainer.GetDSN(), "../../migrations")
	s.Require().NoError(err)

	poolConfig, err := pgxpool.ParseConfig(psqlContainer.GetDSN())
	s.Require().NoError(err)

	poolConfig.MaxConns = int32(cfg.Storage.MaxConnections)
	poolConfig.MinConns = int32(cfg.Storage.MinConnections)
	poolConfig.MaxConnLifetime = time.Duration(cfg.Storage.MaxLifeTime)
	poolConfig.MaxConnIdleTime = time.Duration(cfg.Storage.MaxIdleTime)
	poolConfig.HealthCheckPeriod = time.Duration(cfg.Storage.HealthCheckPeriod)

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	s.Require().NoError(err)

	s.pool = pool

	logger := slog.Default()

	ctxManager := manager.NewCtxManager(pool)
	s.ctxManager = ctxManager

	txManager := manager.NewTxManager(pool, slog.Default())

	userRepo := repository.NewPostgresUserRepository(ctxManager)
	pvzRepo := repository.NewPostgresPvzRepository(ctxManager)
	receptionRepo := repository.NewPostgresReceptionRepository(ctxManager)
	productRepo := repository.NewPostgresProductRepository(ctxManager)

	cities := []string{"Moscow", "Kazan"}

	tokenService := token.NewToken(cfg.SecretKey)
	passwordHasher := hasher.NewHasher()

	authService := service.NewAuthService(userRepo, txManager, tokenService, passwordHasher, logger)
	pvzService := service.NewPVZService(pvzRepo, receptionRepo, cities, productRepo, txManager, logger)

	s.authService = authService
	s.pvzService = pvzService
}

func (s *TestSuite) SetupTest() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	defer db.Close()

	_, err = db.Exec(`
        TRUNCATE TABLE users, pvz, reception, product RESTART IDENTITY CASCADE;
    `)
	s.Require().NoError(err)
}

func (s *TestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.psqlContainer.Terminate(ctx))

	s.server.Close()
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
