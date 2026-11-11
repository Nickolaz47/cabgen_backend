package queue

import (
	"context"
	"errors"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func NewAsynqClient(redisAddr string) (*asynq.Client, error) {
	if redisAddr == "" {
		return nil, errors.New("Redis address is empty")
	}

	opt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		PoolSize: 120,
		DB:       0,
	}

	client := asynq.NewClient(opt)
	if err := client.Ping(); err != nil {
		return nil, err
	}

	return client, nil
}

func NewAsynqInspector(redisAddr string) (*asynq.Inspector, error) {
	if redisAddr == "" {
		return nil, errors.New("Redis address is empty")
	}

	opt := asynq.RedisClientOpt{
		Addr: redisAddr,
		DB:   0,
	}

	return asynq.NewInspector(opt), nil
}

func RunServerWithGracefulShutdown(srv *asynq.Server, mux *asynq.ServeMux,
	logger *zap.Logger) {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	errCh := make(chan error, 1)
	go func() { errCh <- srv.Run(mux) }()

	select {
	case <-ctx.Done():
		logger.Info("Shutdown signal received, gracefully stopping worker...")
		srv.Shutdown()
	case err := <-errCh:
		if err != nil {
			logger.Fatal("Worker execution failed.", zap.Error(err))
		}
		return
	}

	if err := <-errCh; err != nil {
		logger.Fatal("Worker execution failed.", zap.Error(err))
	}
}
