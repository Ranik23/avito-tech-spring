package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)



type Config struct {           
	HTTPServer   	HTTPServerConfig   	`yaml:"HTTPServer"`
	GRPCServer   	GRPCServerConfig   	`yaml:"GRPCServer"`
	GatewayServer 	GatewayServer		`yaml:"GatewayServer"`
	Storage      	StorageConfig      	`yaml:"Storage"`
	MetricServer	MetricServerConfig	`yaml:"MetricServer"`
	SecretKey    	string				`yaml:"-"`
	Cities		 	[]string			`yaml:"Cities"`
}

func LoadConfig(configPath, envPath string) (*Config, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		return nil, err
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(configPath)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viper.AutomaticEnv() // для docker override

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err 
	}

	config.SecretKey = viper.GetString("SECRET_KEY")
	
	config.Storage.Host = viper.GetString("DB_HOST")
	config.Storage.Password = viper.GetString("DB_PASSWORD")
	config.Storage.Name = viper.GetString("DB_NAME")
	config.Storage.Username = viper.GetString("DB_USERNAME")


	return &config, nil
}