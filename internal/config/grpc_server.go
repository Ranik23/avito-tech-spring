package config



type GRPCServerConfig struct {
	Host string			`yaml:"host"`
	Port string			`yaml:"port"`
}