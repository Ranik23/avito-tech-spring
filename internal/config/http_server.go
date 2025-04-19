package config

type HTTPServerConfig struct {
	Host 			string		`yaml:"host"`
	Port 			string		`yaml:"port"`
	ShutdownTimeout int64		`yaml:"ShutDown"`
}
