package cmd

import (
	"context"
	"errors"

	"github.com/hugocortes/tools/src/domain"
	_spinRepo "github.com/hugocortes/tools/src/repos/spin"
	_spinUsecase "github.com/hugocortes/tools/src/usecases/spin"
	"github.com/spf13/cobra"
)

var (
	token   string
	gateURL string
)

func setup(ctx context.Context) (domain.SpinUsecase, error) {
	orcaRepo, err := _spinRepo.New(ctx, gateURL, token)
	if err != nil {
		return nil, err
	}
	return _spinUsecase.New(ctx, orcaRepo)
}

func validateGlobalFlags() error {
	if token == "" {
		return errors.New(("token is required"))
	}
	if gateURL == "" {
		return errors.New(("gate is required"))
	}
	return nil
}

func copy() *cobra.Command {
	var from string
	var to string
	var pipeline string

	validateCmdFlags := func() error {
		if from == "" {
			return errors.New("from is required")
		}
		if to == "" {
			return errors.New("to is required")
		}
		if pipeline == "" {
			return errors.New("pipeline is required")
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:     "copy",
		Aliases: []string{"cp"},
		Short:   "copy a pipeline to another application",
		RunE: func(c *cobra.Command, args []string) error {
			if err := validateGlobalFlags(); err != nil {
				return err
			}
			if err := validateCmdFlags(); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.PersistentFlags().StringVarP(&from, "from", "f", "", "from application")
	cmd.PersistentFlags().StringVarP(&to, "to", "t", "", "to application")
	cmd.PersistentFlags().StringVarP(&from, "pipeline", "p", "", "pipeline to copy")

	return cmd
}

func clean() *cobra.Command {
	var app string

	validateCmdFlags := func() error {
		if app == "" {
			return errors.New("app is required")
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:     "clean",
		Aliases: []string{"c"},
		Short:   "manage hanging executions",
		RunE: func(c *cobra.Command, args []string) error {
			if err := validateGlobalFlags(); err != nil {
				return err
			}
			if err := validateCmdFlags(); err != nil {
				return err
			}

			ctx := context.Background()
			orcaUsecase, err := setup(ctx)
			if err != nil {
				return err
			}
			return orcaUsecase.RemoveNonRunningExecutions(ctx, app)
		},
	}

	cmd.PersistentFlags().StringVarP(&app, "app", "a", "", "application to clean")

	return cmd
}

func spinCmd() *cobra.Command {
	spinCmd := &cobra.Command{
		Use:     "spin",
		Aliases: []string{},
		Short:   "manage spin executions",
	}

	spinCmd.PersistentFlags().StringVarP(&gateURL, "gate", "g", "", "gate url")
	spinCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "oauth2 access token")

	// spinCmd.PersistentFlags().StringVarP(&orcaURL, "orca", "o", "", "orca url")

	cleanCmd := clean()
	copyCmd := copy()

	spinCmd.AddCommand(cleanCmd)
	spinCmd.AddCommand(copyCmd)

	return spinCmd
}
