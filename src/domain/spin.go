package domain

import "context"

const (
	PipelineExecutionSucceeded = "SUCCEEDED"
	PipelineExecutionFailed    = "FAILED"
	PipelineExecutionRunning   = "RUNNING"
)

// PipelineExecution ...
type PipelineExecution struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type SpinUsecase interface {
	RemoveNonRunningExecutions(ctx context.Context, appID string) error
}

type SpinRepo interface {
	GetApplicationPipelineExecutions(ctx context.Context, appID string) ([]*PipelineExecution, error)
	DeletePipelineExecution(ctx context.Context, ID string) error
}
