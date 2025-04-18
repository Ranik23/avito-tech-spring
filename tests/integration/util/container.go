package util

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgreSQLContainer struct {
	testcontainers.Container
	Config ContainerConfig
}

type ContainerConfig struct {
	ImageTag   string
	User       string
	Password   string
	MappedPort string
	Database   string
	Host       string
}


func (c *PostgreSQLContainer) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Config.User, c.Config.Password, c.Config.Host, c.Config.MappedPort, c.Config.Database)
}

func NewPostgreSQLContainer(ctx context.Context) (*PostgreSQLContainer, error) {
	cfg := ContainerConfig{
		User:     "postgres",
		Password: "postgres",
		Database: "avito",
	}

	containerPort := "5432/tcp"
	containerRequest := testcontainers.ContainerRequest{
		Env: map[string]string{
			"POSTGRES_USER":     cfg.User,
			"POSTGRES_PASSWORD": cfg.Password,
			"POSTGRES_DB":       cfg.Database,
		},
		ExposedPorts: []string{containerPort},
		Image:        "postgres:latest",
		WaitingFor:   wait.ForListeningPort(nat.Port(containerPort)),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start PostgreSQL container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(containerPort))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port for container: %w", err)
	}

	cfg.Host = host
	cfg.MappedPort = mappedPort.Port()

	if err := container.Start(context.Background()); err != nil {
		return nil, err
	}

	fmt.Println("PostgreSQL container started on", cfg.Host, cfg.MappedPort)

	return &PostgreSQLContainer{
		Container: container,
		Config:    cfg,
	}, nil
}
