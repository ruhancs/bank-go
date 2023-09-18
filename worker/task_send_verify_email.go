package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

//nome para a fila reconhecer o processo
const TaskSendVerifyEmail = "task:send_verify_email"

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail (
	ctx context.Context, 
	payload *PayloadSendVerifyEmail, 
	opts ...asynq.Option,
) error {
	jsonPayload,err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	//criar a tarefa para a fila
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	//enviar a tarefa para a fila
	info,err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("enqueued task: %v",info)
	return nil
}

func(processor *RedisProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	_,err := processor.store.GetUser(ctx,payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesnt exist: %w", err)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	log.Println("processed task")

	return nil
}