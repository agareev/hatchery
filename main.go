package main

import (
	"github.com/go-openapi/loads"
	"github.com/sirupsen/logrus"
	"hatchery/components"
	"hatchery/handlers"
	"hatchery/server/restapi"
	"hatchery/server/restapi/operations"
)

func main() {

	cfg := components.NewConfiguration()

	//if err := cfg.ParseEnvManual(); err != nil {
	//	logrus.Fatal("Error upload conf from env: ", err)
	//}

	if err := cfg.ParseEnvPkg(); err != nil {
		logrus.Fatal("Error upload conf from env: ", err)
	}

	dbComponent, err := components.NewStorageComponent(cfg)
	if err != nil {
		logrus.Fatal("Error connect to DB", err)
	}

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		logrus.Fatal("Error init swagger spec", err)
	}

	api := operations.NewHatcheryAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Host = cfg.Rest.Host
	server.Port = cfg.Rest.Port

	h := handlers.NewNamesHandler(dbComponent)
	//custom
	server.SetHandler(restapi.ConfigureAPI(api, h))

	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		logrus.Fatal("Error server run", err)
	}
	logrus.Infof("Server run on host %s and port %d", server.Host, server.Port)
}
