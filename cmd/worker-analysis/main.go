package main

import (
	"fmt"
	"log"

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

	// Fail fast: refuse to start without all DB paths configured.
	// In dev, only log — mounts may be empty on purpose.
	inDev := config.Environment == "dev" || config.Environment == ""
	dbPaths := []struct{ name, path string }{
		{"RESFINDER_DB_PATH", config.ResfinderDBPath},
		{"KRAKEN_DB_PATH", config.KrakenDBPath},
		{"BAKTA_DB_PATH", config.BaktaDBPath},
		{"CHECKM_DATA_PATH", config.CheckMDataPath},
		{"FASTANI_PATH", config.FastaniPath},
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

	missing := 0
	for _, db := range dbPaths {
		err := utils.ValidateDBPath(db.path)
		if err == nil {
			continue
		}
		missing++
		msg := fmt.Sprintf("%s: %v", db.name, err)
		if inDev {
			log.Printf("WARNING: %s", msg)
			continue
		}
		log.Fatal(msg)
	}

	if inDev && missing > 0 {
		log.Printf("WARNING: started in dev with %d of %d DB paths invalid; "+
			"analyses using them will fail",
			missing, len(dbPaths))
	}

	// Setup database
	mainDriver := "postgres"
	mainDSN := config.DatabaseConnectionString

	mainDB, err := db.NewGormDatabase(mainDriver, mainDSN)
	if err != nil {
		log.Fatal(err)
	}

	// Logs
	if err := logging.SetupLoggers("./logs/worker-analysis.log"); err != nil {
		log.Fatal(err)
	}
	defer logging.FileLogger.Sync()

	// Asynq Client
	asynqClient, err := queue.NewAsynqClient(config.RedisURL)
	if err != nil {
		log.Fatal(err)
	}

	// One async dependency for the runner (enqueue/delete)
	asynqInspector, err := queue.NewAsynqInspector(config.RedisURL)
	if err != nil {
		log.Fatal(err)
	}
	taskQueue := queue.NewTaskQueue(asynqClient, asynqInspector)

	// Analysis Runner Service
	toolsConfig := pipeline.ToolsConfig{
		FastQCPath:         config.FastQCPath,
		UnicyclerPath:      config.UnicyclerPath,
		SpadesPath:         config.SpadesPath,
		CheckMPath:         config.CheckMPath,
		CheckMDataPath:     config.CheckMDataPath,
		Kraken2Path:        config.Kraken2Path,
		KrakenDBPath:       config.KrakenDBPath,
		BaktaDBPath:        config.BaktaDBPath,
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
		mainDB.DB(), toolsConfig, taskQueue, rootDir, logging.FileLogger,
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
			Concurrency: config.AnalysisConcurrency,
			Queues: map[string]int{
				tasks.QueueAnalysis: 1,
			},
			Logger: logging.FileLogger.Sugar(),
		},
	)

	logging.FileLogger.Info("Starting CABGen Analysis Worker...",
		zap.String("redis_addr", config.RedisURL),
		zap.Int("concurrency", config.AnalysisConcurrency))

	if err := srv.Run(mux); err != nil {
		logging.FileLogger.Fatal("Analysis worker execution failed.",
			zap.Error(err))
	}
}
