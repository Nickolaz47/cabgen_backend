package main

import (
	"context"
	"log"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/container"
	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/routes/admin"
	"github.com/CABGenOrg/cabgen_backend/internal/routes/common"
	"github.com/CABGenOrg/cabgen_backend/internal/routes/public"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load env
	if err := config.LoadEnvVariables(""); err != nil {
		log.Fatal(err)
	}

	// Setup main database
	mainDriver := "postgres"
	mainDSN := config.DatabaseConnectionString
	modelsToMigrate := []any{
		&models.User{},
		&models.Country{},
		&models.Origin{},
		&models.Sequencer{},
		&models.SampleSource{},
		&models.Laboratory{},
		&models.Microorganism{},
	}

	mainDB, err := db.NewGormDatabase(mainDriver, mainDSN)
	if err != nil {
		log.Fatal(err)
	}

	if err := mainDB.Migrate(modelsToMigrate...); err != nil {
		log.Fatal(err)
	}

	// Event Database
	eventDSN := "events.db?_journal_mode=WAL&_busy_timeout=5000"
	eventDB, err := db.NewGormDatabase("sqlite", eventDSN)
	if err != nil {
		log.Fatal(err)
	}

	modelsToMigrate = []any{
		&models.Event{},
	}

	if err := eventDB.Migrate(modelsToMigrate...); err != nil {
		log.Fatal(err)
	}

	eventRepo := container.BuildEventRepository(eventDB.DB())

	// Load translations
	translation.LoadTranslation()

	// Seeder
	if err := utils.Setup(mainDB.DB()); err != nil {
		log.Fatal(err)
	}

	// Logs
	logging.SetupLoggers("./logs/api.log")
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
	authSvc := container.BuildAuthService(mainDB.DB(), eventDB.DB(), logging.FileLogger)
	userSvc := container.BuildUserService(mainDB.DB(), logging.FileLogger)
	admUserSvc := container.BuildAdminUserService(mainDB.DB(), logging.FileLogger)
	labSvc := container.BuildLaboratoryService(mainDB.DB(), logging.FileLogger)
	sequencerSvc := container.BuildSequencerService(mainDB.DB(), logging.FileLogger)
	originSvc := container.BuildOriginService(mainDB.DB(), logging.FileLogger)
	sampleSourceSvc := container.BuildSampleSourceService(mainDB.DB(), logging.FileLogger)
	countrySvc := container.BuildCountryService(mainDB.DB(), logging.FileLogger)
	emailSvc := container.BuildEmailService(mainDB.DB(), logging.FileLogger)

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

	// Event dispatcher
	registry := container.BuildRegistry(emailSvc)
	dispatcher := container.BuildEventDispatcher(eventRepo, registry)

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	go dispatcher.Run(ctx)

	r.Run()
}
