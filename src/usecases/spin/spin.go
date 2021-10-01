package spin

import (
	"context"
	"errors"

	"github.com/hugocortes/tools/src/domain"
	"golang.org/x/sync/errgroup"
)

type usecase struct {
	spinRepo domain.SpinRepo
}

func New(context context.Context, spinRepo domain.SpinRepo) (domain.SpinUsecase, error) {
	return &usecase{
		spinRepo: spinRepo,
	}, nil
}

func (u *usecase) CopyTemplatedPipeline(ctx context.Context, pipeline, from, to string) error {
	fromConfigs, err := u.spinRepo.GetApplicationPipelineConfigs(ctx, from)
	if err != nil {
		return err
	}
	toConfigs, err := u.spinRepo.GetApplicationPipelineConfigs(ctx, to)
	if err != nil {
		return err
	}

	originalConfig := &domain.PipelineConfig{}
	for _, fromConfig := range fromConfigs {
		for _, toConfig := range toConfigs {
			toContainsConfig := fromConfig.Name == toConfig.Name
			if toContainsConfig {
				return errors.New("to application contains pipeline")
			}
		}

		if fromConfig.Name == pipeline {
			originalConfig = fromConfig
		}
	}
	if originalConfig.Name == "" {
		return errors.New("pipeline not found in from application")
	}

	config := &domain.PipelineConfig{
		Schema:               "v2",
		Application:          to,
		Name:                 pipeline,
		Template:             originalConfig.Template,
		Variables:            originalConfig.Variables,
		Type:                 "templatedPipeline",
		LimitConcurrent:      true,
		KeepWaitingPipelines: false,
	}

	_, err = u.spinRepo.CreatePipeline(ctx, config)

	return err
}

func (u *usecase) RemoveNonRunningExecutions(ctx context.Context, application string) error {
	executions, err := u.spinRepo.GetApplicationPipelineExecutions(ctx, application)
	if err != nil {
		return err
	}

	group, ctx := errgroup.WithContext(ctx)
	for _, exe := range executions {
		if exe.ID != "" && exe.Status != domain.PipelineExecutionRunning {
			exeID := exe.ID
			group.Go(func() error {
				return u.spinRepo.DeletePipelineExecution(ctx, exeID)
			})
		}
	}

	err = group.Wait()
	return err
}
