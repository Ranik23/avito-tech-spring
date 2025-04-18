package config

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)


type StorageConfig struct {
	Host               string     `yaml:"host"`
	Port               string     `yaml:"port"`
	Name           	   string     `yaml:"name"`
	Username           string     `yaml:"username"`
	Password           string     `yaml:"password"`
	SSLMode            string     `yaml:"ssl_mode"`


	MaxConnections    int 		  `yaml:"max_connections"`
	MinConnections    int    	  `yaml:"min_connections"`
	MaxLifeTime       int 		  `yaml:"max_lifetime"`
	MaxIdleTime       int 		  `yaml:"max_idle_time"`
	HealthCheckPeriod int 		  `yaml:"health_check_period"`
}

func (s *StorageConfig) Connect() (*pgxpool.Pool, error) {

	dsn := s.GetDSN()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(s.MaxConnections)
	poolConfig.MinConns = int32(s.MinConnections)
	poolConfig.MaxConnLifetime = time.Duration(s.MaxLifeTime)
	poolConfig.MaxConnIdleTime = time.Duration(s.MaxIdleTime)
	poolConfig.HealthCheckPeriod = time.Duration(s.HealthCheckPeriod)

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}


func (s *StorageConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		s.Host, s.Port, s.Username, s.Password, s.Name, s.SSLMode,
	)
}