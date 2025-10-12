package main

import (
	"log"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/routes"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func init() {
	if err := config.LoadEnvVariables(""); err != nil {
		log.Fatal(err)
	}

	driver := "postgres"
	dns := config.DatabaseConnectionString
	if err := db.Connect(driver, dns); err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	repository.InitRepositories(db.DB)
	translation.LoadTranslation()

	if err := utils.Setup(); err != nil {
		log.Fatal(err)
	}

	logging.SetupLoggers("./logs/api.log")
}

func main() {
	defer logging.ConsoleLogger.Sync()
	defer logging.FileLogger.Sync()

	gin.SetMode(gin.DebugMode)
	if config.Environment != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(
		middlewares.LoggerMiddleware(logging.ConsoleLogger, logging.FileLogger),
		middlewares.I18nMiddleware(),
		gin.Recovery(),
	)

	api := r.Group("/api")
	routes.Router(api)

	r.Run()
}
