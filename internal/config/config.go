package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)



type Config struct {           
	HTTPServer   HTTPServerConfig   `yaml:"http"`
	GRPCServer   GRPCServerConfig   `yaml:"grpc"`
	Storage      StorageConfig      `yaml:"storage"`
	SecretKey    string				`yaml:"secret"`
	Cities		 []string			`yaml"cities"`
}

func LoadConfig(configPath, envPath string) (*Config, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		return nil, err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	viper.BindEnv("storage.host", "DB_HOST") //nolint
	viper.BindEnv("storage.password", "DB_PASSWORD") //nolint

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err 
	}

	return &config, nil
}