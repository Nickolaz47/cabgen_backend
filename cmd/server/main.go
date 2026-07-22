package main

import (
	"log"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/container"
	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/queue"
	"github.com/CABGenOrg/cabgen_backend/internal/routes/admin"
	"github.com/CABGenOrg/cabgen_backend/internal/routes/common"
	"github.com/CABGenOrg/cabgen_backend/internal/routes/public"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Root dir
	rootDir, err := utils.GetProjectRoot()
	if err != nil {
		log.Fatal(err)
	}

	// Load env
	if err := config.LoadEnvVariables(""); err != nil {
		log.Fatal(err)
	}

	// Setup database
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
		&models.HealthService{},
		&models.Sample{},
		&models.Analysis{},
		&models.Ticket{},
		&models.PasswordReset{},
	}

	mainDB, err := db.NewGormDatabase(mainDriver, mainDSN)
	if err != nil {
		log.Fatal(err)
	}

	if err := mainDB.Migrate(modelsToMigrate...); err != nil {
		log.Fatal(err)
	}

	// Asynq Client
	asynqClient, err := queue.NewAsynqClient(config.RedisURL)
	if err != nil {
		log.Fatal(err)
	}

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

	// API Config
	r := gin.New()

	corsConfig := cors.Config{}
	corsConfig.AllowOrigins = []string{config.FrontendURL}
	corsConfig.AllowMethods = []string{
		"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD",
	}
	corsConfig.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"Accept-Language",
		"X-Requested-With",
		"Cache-Control"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour

	gin.SetMode(gin.DebugMode)
	if config.Environment != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(
		cors.New(corsConfig),
		middlewares.LoggerMiddleware(logging.ConsoleLogger, logging.FileLogger),
		middlewares.I18nMiddleware(),
		gin.Recovery(),
	)

	api := r.Group("/api")

	// Services
	authSvc := container.BuildAuthService(mainDB.DB(), asynqClient,
		logging.FileLogger)

	userSvc := container.BuildUserService(mainDB.DB(), logging.FileLogger)
	admUserSvc := container.BuildAdminUserService(mainDB.DB(), asynqClient,
		logging.FileLogger)
	labSvc := container.BuildLaboratoryService(mainDB.DB(), logging.FileLogger)
	sequencerSvc := container.BuildSequencerService(mainDB.DB(),
		logging.FileLogger)
	originSvc := container.BuildOriginService(mainDB.DB(), logging.FileLogger)
	sampleSourceSvc := container.BuildSampleSourceService(mainDB.DB(),
		logging.FileLogger)
	countrySvc := container.BuildCountryService(mainDB.DB(), logging.FileLogger)
	microSvc := container.BuildMicroorganismService(mainDB.DB(),
		logging.FileLogger)
	healthServiceSvc := container.BuildHealthServiceService(mainDB.DB(),
		logging.FileLogger)
	sampleSvc := container.BuildSampleService(mainDB.DB(), rootDir,
		logging.FileLogger)
	analysisSvc := container.BuildAnalysisService(mainDB.DB(), asynqClient,
		logging.FileLogger)
	admAnalysisSvc := container.BuildAdminAnalysisService(mainDB.DB(),
		asynqClient, logging.FileLogger)
	ticketSvc := container.BuildTicketService(mainDB.DB(), asynqClient,
		logging.FileLogger)

	// Public handlers
	healthHandler := container.BuildHealthHandler()
	pubAuthHandler := container.BuildPublicAuthHandler(authSvc)
	pubCountryHandler := container.BuildPublicCountryHandler(countrySvc)
	contactHandler := container.BuildTicketHandler(ticketSvc)

	// Common handlers
	authHandler := container.BuildCommonAuthHandler(authSvc)
	userHandler := container.BuildUserHandler(userSvc)
	laboratoryHandler := container.BuildLaboratoryHandler(labSvc)
	sequencerHandler := container.BuildSequencerHandler(sequencerSvc)
	originHandler := container.BuildOriginHandler(originSvc)
	sampleSourceHandler := container.BuildSampleSourceHandler(sampleSourceSvc)
	microHandler := container.BuildMicroorganismHandler(microSvc)
	healthServiceHandler := container.BuildHealthServiceHandler(healthServiceSvc)
	sampleHandler := container.BuildSampleHandler(sampleSvc)
	analysisHandler := container.BuildAnalysisHandler(analysisSvc)
	selectOptionHandler := container.BuildSelectOptionHandler()
	cityHandler := container.BuildCityHandler()

	// Admin handlers
	adminUserHandler := container.BuildAdminUserHandler(admUserSvc)
	adminLaboratoryHandler := container.BuildAdminLaboratoryHandler(labSvc)
	adminSequencerHandler := container.BuildAdminSequencerHandler(sequencerSvc)
	adminOriginHandler := container.BuildAdminOriginHandler(originSvc)
	adminSampleSourceHandler := container.BuildAdminSampleSourceHandler(
		sampleSourceSvc)
	adminCountryHandler := container.BuildAdminCountryHandler(countrySvc)
	adminMicroHandler := container.BuildAdminMicroorganismHandler(microSvc)
	adminHealthServiceHandler := container.BuildAdminHealthServiceHandler(
		healthServiceSvc)
	adminSampleHandler := container.BuildAdminSampleHandler(sampleSvc)
	adminAnalysisHandler := container.BuildAdminAnalysisHandler(admAnalysisSvc)
	adminTicketHandler := container.BuildAdminTicketHandler(ticketSvc)

	// Public routes
	publicRouter := api.Group("")
	public.SetupCountryRoutes(publicRouter, pubCountryHandler)
	public.SetupHealthRoute(publicRouter, healthHandler)
	public.SetupPublicAuthRoutes(publicRouter, pubAuthHandler)
	public.SetupContactRoutes(publicRouter, contactHandler)

	// Common routes
	commonRouter := api.Group("", middlewares.AuthMiddleware())
	common.SetupCommonAuthRoutes(commonRouter, authHandler)
	common.SetupUserRoutes(commonRouter, userHandler)
	common.SetupSequencerRoutes(commonRouter, sequencerHandler)
	common.SetupLaboratoryRoutes(commonRouter, laboratoryHandler)
	common.SetupOriginRoutes(commonRouter, originHandler)
	common.SetupSampleSourceRoutes(commonRouter, sampleSourceHandler)
	common.SetupMicroorganismRoutes(commonRouter, microHandler)
	common.SetupHealthServiceRoutes(commonRouter, healthServiceHandler)
	common.SetupSampleRoutes(commonRouter, sampleHandler)
	common.SetupAnalysisRoutes(commonRouter, analysisHandler)
	common.SetupSelectOptionRoutes(commonRouter, selectOptionHandler)
	common.SetupCityRoutes(commonRouter, cityHandler)

	// Admin routes
	adminRouter := api.Group("/admin", middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware())
	admin.SetupAdminUserRoutes(adminRouter, adminUserHandler)
	admin.SetupAdminSequencerRoutes(adminRouter, adminSequencerHandler)
	admin.SetupAdminLaboratoryRoutes(adminRouter, adminLaboratoryHandler)
	admin.SetupAdminOriginRoutes(adminRouter, adminOriginHandler)
	admin.SetupAdminSampleSourceRoutes(adminRouter, adminSampleSourceHandler)
	admin.SetupAdminCountryRoutes(adminRouter, adminCountryHandler)
	admin.SetupAdminMicroorganismRoutes(adminRouter, adminMicroHandler)
	admin.SetupAdminHealthServiceRoutes(adminRouter, adminHealthServiceHandler)
	admin.SetupAdminSampleRoutes(adminRouter, adminSampleHandler)
	admin.SetupAdminAnalysisRoutes(adminRouter, adminAnalysisHandler)
	admin.SetupAdminTicketRoutes(adminRouter, adminTicketHandler)

	r.Run()
}
