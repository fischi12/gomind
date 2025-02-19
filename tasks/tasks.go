package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"gomind/database"
	"gomind/services"
	"gorm.io/gorm"
)

const (
	TypeHandAbstractionFlop = "hand_abstraction:flop"
)

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

func NewHandAbstractionFlopTask(hand []services.Hand) (*asynq.Task, error) {
	payload, err := json.Marshal(hand)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeHandAbstractionFlop, payload), nil
}

// ---------------------------------------------------------------
// Write a function HandleXXXTask to handle the input task.
// Note that it satisfies the asynq.HandlerFunc interface.
//
// Handler doesn't need to be a function. You can define a type
// that satisfies asynq.Handler interface. See examples below.
// ---------------------------------------------------------------
type TaskHandler struct {
	DB *gorm.DB
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		DB: database.GetDB(),
	}
}

func (h *TaskHandler) HandleHandAbstractionFlopTask(ctx context.Context, t *asynq.Task) error {
	var hand []services.Hand
	if err := json.Unmarshal(t.Payload(), &hand); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	return services.CalculateAndSaveHandStrength(hand, h.DB)
}
