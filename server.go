//usr/bin/env go run $0 $@; exit;
package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/aufaitio/listener/apis"
	"github.com/aufaitio/listener/app"
	"github.com/aufaitio/listener/daos"
	"github.com/aufaitio/listener/errors"
	"github.com/aufaitio/listener/services"
	"github.com/docopt/docopt-go"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/cors"
)

func main() {
	doc := `Au Fait
Command line interface for starting the listener API for Au Fait.

Usage:
	server [--configPath=<path>]
Options:
	-h --help				Show this message
	--version				Show version info
	--configPath=<path>   	Path to app.yaml config file [default: config]`

	arguments, _ := docopt.ParseDoc(doc)
	configPath := arguments["--configPath"].(string)
	fmt.Println(configPath)

	// load application configurations
	if err := app.LoadConfig(configPath); err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}

	// load error messages
	if err := errors.LoadMessages(app.Config.ErrorFile); err != nil {
		panic(fmt.Errorf("Failed to read the error message file: %s", err))
	}

	// create the logger
	logger := logrus.New()

	// connect to the database
	client, err := mongo.Connect(context.Background(), buildDBHost(app.Config), nil)

	if err != nil {
		panic(fmt.Errorf("Failed to connect to MongoDB with error message: %s", err))
	}

	db := client.Database(app.Config.DB.Name)

	// wire up API routing
	http.Handle("/", buildRouter(logger, db))

	// start the server
	address := fmt.Sprintf(":%v", app.Config.Port)
	logger.Infof("server %v is started at %v\n", app.Version, address)
	panic(http.ListenAndServe(address, nil))
}

func buildDBHost(config app.AppConfig) string {
	prefix := ""

	if config.DB.Username != "" {
		prefix = fmt.Sprintf("%s:%s@", config.DB.Username, config.DB.Password)
	}

	return fmt.Sprintf("mongodb://%s%s:%d", prefix, config.DB.Host, config.DB.Port)
}

func buildRouter(logger *logrus.Logger, db *mongo.Database) *routing.Router {
	router := routing.New()

	router.To("GET,HEAD", "/heartbeat", func(c *routing.Context) error {
		c.Abort() // skip all other middlewares/handlers
		return c.Write("OK " + app.Version)
	})

	router.Use(
		app.Init(logger, db),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.Options{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "*",
		}),
	)

	rg := router.Group("/v1")

	repoDAO := daos.NewRepositoryDAO()
	repoService := services.NewRepositoryService(repoDAO)
	apis.ServeRepositoryResource(rg, repoService)
	jobDAO := daos.NewJobDAO()
	apis.ServeJobResource(rg, services.NewJobService(jobDAO, repoDAO), repoService)

	return router
}
