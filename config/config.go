package config

import "os"

// Config is the struct that holds env variables
type Config struct {
	DockerHost      string
	MongoHost       string
	MongoName       string
	MongoUser       string
	MongoPass       string
	MongoCollection string
}

// SetConfigs set all needed env variables
// retornar erro!
func (c Config) SetConfigs() {
	// check https://github.com/spf13/viper
	c.DockerHost = os.Getenv("DOCKER_HOST")
	c.MongoHost = os.Getenv("MONGO_HOST")
	c.MongoName = os.Getenv("MONGO_NAME")
	c.MongoUser = os.Getenv("MONGO_USER")
	c.MongoPass = os.Getenv("MONGO_PASS")
	c.MongoCollection = os.Getenv("MONGO_COLLECTION")
}
