package config

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)


type StorageConfig struct {
	Host               string     `yaml:"Host"`
	Port               string     `yaml:"Port"`
	Name           	   string     `yaml:"-"`
	Username           string     `yaml:"-"`
	Password           string     `yaml:"-"`
	SSLMode            string     `yaml:"ssl"`


	MaxConnections    int 		  `yaml:"MaxConnections"`
	MinConnections    int    	  `yaml:"MinConnections"`
	MaxLifeTime       int 		  `yaml:"MaxLifetime"`
	MaxIdleTime       int 		  `yaml:"MaxIdleTime"`
	HealthCheckPeriod int 		  `yaml:"HealthCheckPeriod"`
}

func (s *StorageConfig) Connect() (*pgxpool.Pool, error) {

	dsn := s.GetDSN()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(s.MaxConnections)
	poolConfig.MinConns = int32(s.MinConnections)
	poolConfig.MaxConnLifetime = time.Duration(s.MaxLifeTime) * time.Second
	poolConfig.MaxConnIdleTime = time.Duration(s.MaxIdleTime) * time.Second
	poolConfig.HealthCheckPeriod = time.Duration(s.HealthCheckPeriod) * time.Second

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