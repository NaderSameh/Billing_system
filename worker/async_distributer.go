package worker

import "github.com/hibiken/asynq"

type TaskDistributer interface {
	NewMonthlyBillingTask(ClientID int64, OrderID int64) (*asynq.Task, error)
}

type RedisTaskDistributer struct {
	client *asynq.Client
}

func NewTaskDistributer(addr string) TaskDistributer {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: addr})
	return &RedisTaskDistributer{
		client: client,
	}
}
