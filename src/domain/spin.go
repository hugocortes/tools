package domain

import (
	"context"
)

const (
	PipelineExecutionSucceeded = "SUCCEEDED"
	PipelineExecutionFailed    = "FAILED"
	PipelineExecutionRunning   = "RUNNING"
)

type PipelineTemplate struct {
	ArtifactAccount string `json:"artifactAccount"`
	Reference       string `json:"reference"`
	Type            string `json:"type"`
}

type PipelineConfig struct {
	Schema               string                 `json:"schema"`
	Application          string                 `json:"application"`
	Name                 string                 `json:"name"`
	Template             PipelineTemplate       `json:"template"`
	Variables            map[string]interface{} `json:"variables"`
	Type                 string                 `json:"type"`
	LimitConcurrent      bool                   `json:"limitConcurrent"`
	KeepWaitingPipelines bool                   `json:"keepWaitingPipelines"`
}

type PipelineExecution struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type SpinUsecase interface {
	RemoveNonRunningExecutions(ctx context.Context, application string) error
	CopyTemplatedPipeline(ctx context.Context, pipeline string, from string, to string) error
}

type SpinRepo interface {
	CreatePipeline(ctx context.Context, pipelineConfig *PipelineConfig) (*PipelineConfig, error)
	GetApplicationPipelineConfigs(ctx context.Context, application string) ([]*PipelineConfig, error)
	GetApplicationPipelineExecutions(ctx context.Context, application string) ([]*PipelineExecution, error)
	DeletePipelineExecution(ctx context.Context, ID string) error
}
