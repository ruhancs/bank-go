package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

//distribuir os tarefas para o redis

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail (
		ctx context.Context, 
		payload *PayloadSendVerifyEmail, 
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}