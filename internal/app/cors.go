package app

import "github.com/gin-contrib/cors"



func NewCORSConfig() cors.Config {
	
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Authorization", "Content-Type"}

	return config
}