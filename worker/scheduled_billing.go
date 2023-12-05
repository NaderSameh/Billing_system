package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeMonthlyBilling = "bill:monthly"
)

type MonthlyBillingPayload struct {
	ClientID int64
	OrderID  int64
}

func (distributor *RedisTaskDistributer) NewMonthlyBillingTask(ClientID int64, OrderID int64) (*asynq.Task, error) {
	payload, err := json.Marshal(MonthlyBillingPayload{ClientID: ClientID, OrderID: OrderID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeMonthlyBilling, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}

func HandleMonthlyBillingTask(ctx context.Context, t *asynq.Task) error {
	var p MonthlyBillingPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	// log.Printf("Sending Email to User: user_id=%d, template_id=%s", p.UserID, p.TemplateID)
	// Email delivery code ...
	return nil
}
