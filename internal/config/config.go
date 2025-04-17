package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)



type Config struct {

}



func LoadConfig(configPath, envPath string) (*Config, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("WARNING: error loading .env file from %s: %v\n", envPath, err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	viper.BindEnv("storage.postgres.hosts", "DB_HOST")
	viper.BindEnv("storage.postgres.password", "DB_PASSWORD")
	viper.BindEnv("jwt.secret_key", "JWT_SECRET_KEY")

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err 
	}

	return &config, nil
}