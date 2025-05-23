//go:build integration

package integration

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/Ranik23/avito-tech-spring/internal/config"
	"github.com/Ranik23/avito-tech-spring/internal/hasher"
	"github.com/Ranik23/avito-tech-spring/internal/repository"
	"github.com/Ranik23/avito-tech-spring/internal/repository/postgresql"
	"github.com/Ranik23/avito-tech-spring/internal/service"
	"github.com/Ranik23/avito-tech-spring/internal/token"
	"github.com/Ranik23/avito-tech-spring/tests/integration/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
	"github.com/stretchr/testify/suite"
)


type TestSuite struct {
	suite.Suite
	psqlContainer *util.PostgreSQLContainer


	authService service.AuthService
	pvzService service.PVZService

	receptionRepo repository.ReceptionRepository
	pvzRepo repository.PvzRepository
	userRepo repository.UserRepository
	productRepo repository.ProductRepository

	token token.Token
	hasher hasher.Hasher

	cities []string

	logger *slog.Logger

	service     service.Service	 
}

func (s *TestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer ctxCancel()

	logger := slog.New(tint.NewHandler(os.Stdout, nil))

	cfg, err := config.LoadConfig("../../config", "../../.env",)
	s.Require().NoError(err)

	psqlContainer, err := util.NewPostgreSQLContainer(ctx)
	s.Require().NoError(err)

	s.psqlContainer = psqlContainer

	err = util.RunMigrations(psqlContainer.GetDSN(), "../../migrations/goose")
	s.Require().NoError(err)

	poolConfig, err := pgxpool.ParseConfig(psqlContainer.GetDSN())
	s.Require().NoError(err)

	poolConfig.MaxConns = int32(cfg.Storage.MaxConnections)
	poolConfig.MinConns = int32(cfg.Storage.MinConnections)
	poolConfig.MaxConnLifetime = time.Duration(cfg.Storage.MaxLifeTime) * time.Second
	poolConfig.MaxConnIdleTime = time.Duration(cfg.Storage.MaxIdleTime) * time.Second
	poolConfig.HealthCheckPeriod = time.Duration(cfg.Storage.HealthCheckPeriod) * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	s.Require().NoError(err)


	ctxManager := postgresql.NewCtxManager(pool)
	txManager := postgresql.NewTxManager(pool, logger, ctxManager)

	userRepo := postgresql.NewPostgresUserRepository(ctxManager, logger)
	productRepo := postgresql.NewPostgresProductRepository(ctxManager, logger)
	receptionRepo := postgresql.NewPostgresReceptionRepository(ctxManager, logger)
	pvzRepo := postgresql.NewPostgresPvzRepository(ctxManager, logger)


	s.pvzRepo = pvzRepo
	s.userRepo = userRepo
	s.productRepo = productRepo
	s.receptionRepo = receptionRepo

	token := token.NewToken("lol", logger)
	hasher := hasher.NewHasher()


	s.token = token
	s.hasher = hasher

	cities := []string{"Moscow"}
	s.cities = cities

	authService := service.NewAuthService(userRepo, txManager, token, hasher, logger)
	pvzService := service.NewPVZService(pvzRepo, receptionRepo, cities, productRepo, txManager, logger)

	service := service.NewService(authService, pvzService)

	s.logger = logger

	s.authService = authService
	s.pvzService = pvzService
	s.service = service
}

func (s *TestSuite) SetupTest() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	defer db.Close()

	_, err = db.Exec(`
        TRUNCATE TABLE users, reception, product, pvz RESTART IDENTITY CASCADE;
    `)
	s.Require().NoError(err)
}

func (s *TestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}



