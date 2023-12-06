package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	TypeEmailDelivery = "email:deliver"
)

type PayloadSendEmail struct {
	OrderID      int64  `json:"order_id"`
	BatchName    string `json:"batch_name"`
	CustomerName string `json:"customer_name"`
	BatchID      int64  `json:"batch_id"`
}

func (distributor *RedisTaskDistributor) NewEmailDeliveryTask(payload *PayloadSendEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TypeEmailDelivery, jsonPayload, opts...)
	info, err := distributor.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	subject := fmt.Sprintf("New order activation %d", payload.OrderID)

	content := fmt.Sprintf(`Hello Cypod engineer,<br/>
	A new order has been activated for batch: %s with ID: %d for customer: %s`, payload.BatchName, payload.BatchID, payload.CustomerName)

	to := []string{viper.GetString("SUPPORT_EMAIL")}

	err := processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Msg("processed task")
		// Str("email", user.Email).Msg("processed task")
	return nil
}
