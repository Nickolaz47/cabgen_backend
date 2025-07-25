package main

import (
	"log"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	err := config.LoadEnvVariables("")
	if err != nil {
		log.Fatal(err)
	}

	db.Connect()
	db.Migrate()
}

func main() {
	r := gin.Default()
	routes.Router(r)

	r.Run()
}
