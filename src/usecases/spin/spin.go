package spin

import (
	"context"

	"github.com/hugocortes/tools/src/domain"
	"golang.org/x/sync/errgroup"
)

type usecase struct {
	orcaRepo domain.SpinRepo
}

func New(context context.Context, orcaRepo domain.SpinRepo) (domain.SpinUsecase, error) {
	return &usecase{
		orcaRepo: orcaRepo,
	}, nil
}

func (u *usecase) RemoveNonRunningExecutions(ctx context.Context, appID string) error {
	executions, err := u.orcaRepo.GetApplicationPipelineExecutions(ctx, appID)
	if err != nil {
		return err
	}

	group, ctx := errgroup.WithContext(ctx)
	for _, exe := range executions {
		if exe.ID != "" && exe.Status != domain.PipelineExecutionRunning {
			exeID := exe.ID
			group.Go(func() error {
				return u.orcaRepo.DeletePipelineExecution(ctx, exeID)
			})
		}
	}

	err = group.Wait()
	return err
}
