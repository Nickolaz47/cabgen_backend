package main

import (
	"log"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/container"
	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/workers"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func main() {
	// Load env
	if err := config.LoadEnvVariables(""); err != nil {
		log.Fatal(err)
	}

	// Setup database
	mainDriver := "postgres"
	mainDSN := config.DatabaseConnectionString

	mainDB, err := db.NewGormDatabase(mainDriver, mainDSN)
	if err != nil {
		log.Fatal(err)
	}

	// Logs
	logging.SetupLoggers("./logs/worker-email.log")
	defer logging.FileLogger.Sync()

	// Email Service
	emailSvc := container.BuildEmailService(mainDB.DB(), logging.FileLogger)

	// Handler
	emailHandler := workers.NewEmailTaskHandler(emailSvc)

	// Mux
	mux := asynq.NewServeMux()
	mux.Handle(tasks.TaskTypeAdminAlertEmail, emailHandler)
	mux.Handle(tasks.TaskTypeWelcomeEmail, emailHandler)
	mux.Handle(tasks.TaskTypeAnalysisDoneEmail, emailHandler)
	mux.Handle(tasks.TaskTypeAdminTicketEmail, emailHandler)
	mux.Handle(tasks.TaskTypeFinishedTicketEmail, emailHandler)
	mux.Handle(tasks.TaskTypePasswordResetEmail, emailHandler)

	// Redis
	redisOpt := asynq.RedisClientOpt{Addr: config.RedisURL}
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				tasks.QueueEmail: 1,
			},
			Logger: logging.FileLogger.Sugar(),
		},
	)

	logging.FileLogger.Info("Starting CABGen Email Worker...",
		zap.String("redis_addr", config.RedisURL),
		zap.Int("concurrency", 10),
	)

	if err := srv.Run(mux); err != nil {
		logging.FileLogger.Fatal("Email worker execution failed.", zap.Error(err))
	}
}
