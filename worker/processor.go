package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/ruhancs/bank-go/db/sqlc"
)

//pegar as tarefas da fila no redis e processar

const (
	QueueCritical = "critical"
	QueueDefalut = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisProcessor struct {
	server *asynq.Server
	store db.Store
}

func NewRedisProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,//prioridade da fila
				QueueDefalut: 5,
			},
		},
	)

	return &RedisProcessor{
		server: server,
		store: store,
	}
}

func (processor *RedisProcessor) Start() error {
	mux := asynq.NewServeMux()

	//similar a criar rota de url so que com tarefas e nomes na fila
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}