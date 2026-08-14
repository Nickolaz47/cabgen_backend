package main

import (
	"log"
	"os"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/container"
	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/pipeline"
	"github.com/CABGenOrg/cabgen_backend/internal/queue"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/workers"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
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

	// Fail fast: refuse to start without all DB paths configured
	dbPaths := []struct{ name, path string }{
		{"RESFINDER_DB_PATH", config.ResfinderDBPath},
		{"KRAKEN_DB_PATH", config.KrakenDBPath},
		{"CHECKM_DATA_PATH", os.Getenv("CHECKM_DATA_PATH")},
		{"POLI_DB_PSEUDO", config.PoliDbPseudo},
		{"POLI_DB_KLEB", config.PoliDbKleb},
		{"POLI_DB_ENTERO", config.PoliDbEntero},
		{"POLI_DB_ACINETO", config.PoliDbAcineto},
		{"OTHER_DB_PSEUDO", config.OtherDbPseudo},
		{"OTHER_DB_KLEB", config.OtherDbKleb},
		{"OTHER_DB_ENTERO", config.OtherDbEntero},
		{"OTHER_DB_ACINETO", config.OtherDbAcineto},
		{"FASTANI_LIST_KLEB", config.FastaniListKleb},
		{"FASTANI_LIST_ENTERO", config.FastaniListEntero},
		{"FASTANI_LIST_ACINETO", config.FastaniListAcineto},
	}
	for _, db := range dbPaths {
		if db.path == "" {
			log.Fatalf("%s is not configured; refusing to start", db.name)
		}
		if _, err := os.Stat(db.path); err != nil {
			log.Fatalf("%s: path %s not accessible: %v", db.name, db.path, err)
		}
	}

	// Setup database
	mainDriver := "postgres"
	mainDSN := config.DatabaseConnectionString

	mainDB, err := db.NewGormDatabase(mainDriver, mainDSN)
	if err != nil {
		log.Fatal(err)
	}

	// Logs
	logging.SetupLoggers("./logs/worker-analysis.log")
	defer logging.FileLogger.Sync()

	// Asynq Client
	asynqClient, err := queue.NewAsynqClient(config.RedisURL)
	if err != nil {
		log.Fatal(err)
	}

	// Analysis Runner Service
	toolsConfig := pipeline.ToolsConfig{
		FastQCPath:         config.FastQCPath,
		UnicyclerPath:      config.UnicyclerPath,
		SpadesPath:         config.SpadesPath,
		CheckMPath:         config.CheckMPath,
		Kraken2Path:        config.Kraken2Path,
		KrakenDBPath:       config.KrakenDBPath,
		FastANIPath:        config.FastaniPath,
		AbricatePath:       config.AbricatePath,
		MLSTPath:           config.MlstPath,
		ResfinderDBPath:    config.ResfinderDBPath,
		PoliDbPseudo:       config.PoliDbPseudo,
		PoliDbKleb:         config.PoliDbKleb,
		PoliDbEntero:       config.PoliDbEntero,
		PoliDbAcineto:      config.PoliDbAcineto,
		OtherDbPseudo:      config.OtherDbPseudo,
		OtherDbKleb:        config.OtherDbKleb,
		OtherDbEntero:      config.OtherDbEntero,
		OtherDbAcineto:     config.OtherDbAcineto,
		FastaniListKleb:    config.FastaniListKleb,
		FastaniListEntero:  config.FastaniListEntero,
		FastaniListAcineto: config.FastaniListAcineto,
	}
	analysisRunnerSvc := container.BuildAnalysisRunnerService(
		mainDB.DB(), toolsConfig, asynqClient, rootDir, logging.FileLogger,
	)

	// Handler
	analysisHandler := workers.NewAnalysisTaskHandler(analysisRunnerSvc,
		logging.FileLogger)

	// Mux
	mux := asynq.NewServeMux()
	mux.Handle(tasks.TaskTypeAnalysisProcess, analysisHandler)

	// Redis
	redisOpt := asynq.RedisClientOpt{Addr: config.RedisURL}
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 4,
			Queues: map[string]int{
				tasks.QueueAnalysis: 1,
			},
			Logger: logging.FileLogger.Sugar(),
		},
	)

	logging.FileLogger.Info("Starting CABGen Analysis Worker...",
		zap.String("redis_addr", config.RedisURL),
		zap.Int("concurrency", 4))

	if err := srv.Run(mux); err != nil {
		logging.FileLogger.Fatal("Analysis worker execution failed.",
			zap.Error(err))
	}
}
