package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config is the server configuration structure.
// all fields will be filled with environment variables.
type Config struct {
	ServerHost string // address that server will listening on
	MongoURL   string // mongo db username
}

// initialize will read environment variables and save them in config structure fields
func (config *Config) initialize() {
	// read environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(0)
	}

	config.ServerHost = os.Getenv("server_host")
	config.MongoURL = os.Getenv("mongo_url")
}

// MongoURI will generate mongo db connect uri
func (config *Config) MongoURI() string {
	return config.MongoURL
}

// NewConfig will create and initialize config struct
func NewConfig() *Config {
	config := new(Config)
	config.initialize()
	return config
}
