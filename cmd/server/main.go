package main

import (
	"log"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/container"
	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/CABGenOrg/cabgen_backend/internal/routes/admin"
	"github.com/CABGenOrg/cabgen_backend/internal/routes/common"
	"github.com/CABGenOrg/cabgen_backend/internal/routes/public"
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

	// Services
	authSvc := container.BuildAuthService(db.DB)
	userSvc := container.BuildUserService(db.DB)
	admUserSvc := container.BuildAdminUserService(db.DB)
	labSvc := container.BuildLaboratoryService(db.DB)
	sequencerSvc := container.BuildSequencerService(db.DB)
	originSvc := container.BuildOriginService(db.DB)
	sampleSourceSvc := container.BuildSampleSourceService(db.DB)
	countrySvc := container.BuildCountryService(db.DB)

	// Public handlers
	healthHandler := container.BuildHealthHandler()
	authHandler := container.BuildAuthHandler(authSvc)
	pubCountryHandler := container.BuildPublicCountryHandler(countrySvc)

	// Common handlers
	userHandler := container.BuildUserHandler(userSvc)
	laboratoryHandler := container.BuildLaboratoryHandler(labSvc)
	sequencerHandler := container.BuildSequencerHandler(sequencerSvc)
	originHandler := container.BuildOriginHandler(originSvc)
	sampleSourceHandler := container.BuildSampleSourceHandler(sampleSourceSvc)

	// Admin handlers
	adminUserHandler := container.BuildAdminUserHandler(admUserSvc)
	adminLaboratoryHandler := container.BuildAdminLaboratoryHandler(labSvc)
	adminSequencerHandler := container.BuildAdminSequencerHandler(sequencerSvc)
	adminOriginHandler := container.BuildAdminOriginHandler(originSvc)
	adminSampleSourceHandler := container.BuildAdminSampleSourceHandler(sampleSourceSvc)
	adminCountryHandler := container.BuildAdminCountryHandler(countrySvc)

	// Public routes
	publicRouter := api.Group("")
	public.SetupCountryRoutes(publicRouter, pubCountryHandler)
	public.SetupHealthRoute(publicRouter, healthHandler)
	public.SetupAuthRoutes(publicRouter, authHandler)

	// Common routes
	commonRouter := api.Group("", middlewares.AuthMiddleware())
	common.SetupUserRoutes(commonRouter, userHandler)
	common.SetupSequencerRoutes(commonRouter, sequencerHandler)
	common.SetupLaboratoryRoutes(commonRouter, laboratoryHandler)
	common.SetupOriginRoutes(commonRouter, originHandler)
	common.SetupSampleSourceRoutes(commonRouter, sampleSourceHandler)

	// Admin routes
	adminRouter := api.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	admin.SetupAdminUserRoutes(adminRouter, adminUserHandler)
	admin.SetupAdminSequencerRoutes(adminRouter, adminSequencerHandler)
	admin.SetupAdminLaboratoryRoutes(adminRouter, adminLaboratoryHandler)
	admin.SetupAdminOriginRoutes(adminRouter, adminOriginHandler)
	admin.SetupAdminSampleSourceRoutes(adminRouter, adminSampleSourceHandler)
	admin.SetupAdminCountryRoutes(adminRouter, adminCountryHandler)

	r.Run()
}
