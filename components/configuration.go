package components

import (
	"github.com/kelseyhightower/envconfig"
	"os"
	"strconv"
)

const (
	appName                 = "myapp"
	dbDefaultUri            = "localhost"
	dbDefaultPort           = "27017"
	dbDefaultName           = "runners_info"
	dbDefaultCollectionName = "numbers"
	dbDefaultHttpHost       = "localhost"
	dbDefaultHttpPort       = 8182
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
			Uri:            dbDefaultUri,
			Port:           dbDefaultPort,
			Name:           dbDefaultName,
			CollectionName: dbDefaultCollectionName,
		},
		Rest: ServerConfig{
			Host: dbDefaultHttpHost,
			Port: dbDefaultHttpPort,
		},
	}
}

func (c *Configuration) ParseEnvManual() (err error) {
	dbUri := os.Getenv(appName + "_DB_URI")
	if dbUri != "" {
		c.Mongo.Uri = dbUri
	}

	dbPort := os.Getenv(appName + "_DB_PORT")
	if dbUri != "" {
		c.Mongo.Port = dbPort
	}

	dbName := os.Getenv(appName + "_DB_NAME")
	if dbName != "" {
		c.Mongo.Name = dbName
	}

	dbCollectionName := os.Getenv(appName + "_DB_COLLECTIONNAME")
	if dbUri != "" {
		c.Mongo.CollectionName = dbCollectionName
	}

	httpHost := os.Getenv(appName + "_HTTP_HOST")
	if httpHost != "" {
		c.Rest.Host = httpHost
	}

	httpPort := os.Getenv(appName + "_HTTP_PORT")
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

	err = envconfig.Process(appName+"_DB", &c.Mongo)
	if err != nil {
		return err
	}

	err = envconfig.Process(appName+"_HTTP", &c.Rest)
	if err != nil {
		return err
	}

	return envconfig.Process("myApp", c)
}
