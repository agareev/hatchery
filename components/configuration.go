package components

import (
	"github.com/kelseyhightower/envconfig"
	"os"
	"strconv"
)

const (
	AppName           = "myapp"
	DbDefaultUri      = "localhost"
	DbDefaultPort     = "27017"
	DbDefaultName     = "runners_info"
	DbDefaultHttpHost = "localhost"
	DbDefaultHttpPort = 8182
)

type Configuration struct {
	Mongo DbConfig
	Rest  ServerConfig
}

type DbConfig struct {
	Uri            string
	Port           string
	Name           string
	CollectionName string
}

type ServerConfig struct {
	Host string
	Port int
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Mongo: DbConfig{
			Uri:  DbDefaultUri,
			Port: DbDefaultPort,
			Name: DbDefaultName,
		},
		Rest: ServerConfig{
			Host: DbDefaultHttpHost,
			Port: DbDefaultHttpPort,
		},
	}
}

func (c *Configuration) ParseEnvManual() (err error) {
	DbUri := os.Getenv(AppName + "_DB_URI")
	if DbUri != "" {
		c.Mongo.Uri = DbUri
	}

	DbPort := os.Getenv(AppName + "_DB_PORT")
	if DbUri != "" {
		c.Mongo.Port = DbPort
	}

	DbName := os.Getenv(AppName + "_DB_NAME")
	if DbName != "" {
		c.Mongo.Name = DbName
	}

	DbCollectionName := os.Getenv(AppName + "_DB_COLLECTIONNAME")
	if DbUri != "" {
		c.Mongo.CollectionName = DbCollectionName
	}

	httpHost := os.Getenv(AppName + "_HTTP_HOST")
	if httpHost != "" {
		c.Rest.Host = httpHost
	}

	httpPort := os.Getenv(AppName + "_HTTP_PORT")
	if httpPort != "" {
		port, err := strconv.Atoi(httpPort)
		if err == nil {
			c.Rest.Port = port
		}
	}

	return nil
}

func (c *Configuration) ParseEnvPkg() (err error) {

	//Look https://github.com/kelseyhightower/envconfig

	err = envconfig.Process(AppName+"_DB", &c.Mongo)
	if err != nil {
		return err
	}

	err = envconfig.Process(AppName+"_HTTP", &c.Rest)
	if err != nil {
		return err
	}

	return envconfig.Process("myApp", c)
}
