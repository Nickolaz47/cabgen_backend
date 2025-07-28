package main

import (
	"log"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/CABGenOrg/cabgen_backend/internal/routes"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func init() {
	if err := config.LoadEnvVariables(""); err != nil {
		log.Fatal(err)
	}

	db.Connect()
	db.Migrate()
	translation.LoadTranslation()

	if err := utils.Setup(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := gin.Default()
	r.Use(middlewares.I18nMiddleware())
	routes.Router(r)

	r.Run()
}
