package cmd

import (
	"context"
	"errors"

	"github.com/hugocortes/tools/spin"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	// flags
	orcaURL string
	app     string
)

func validatePersistentFlags() error {
	if orcaURL == "" {
		return errors.New("orca is required")
	}
	if app == "" {
		return errors.New("app is required")
	}
	return nil
}

func clean() *cobra.Command {
	cleanCmd := &cobra.Command{
		Use:     "clean",
		Aliases: []string{"c"},
		Short:   "manage hanging executions",
		RunE: func(c *cobra.Command, args []string) error {
			if err := validatePersistentFlags(); err != nil {
				return err
			}

			orca, err := spin.NewOrca(orcaURL)
			if err != nil {
				return err
			}
			executions, err := orca.GetExecutions(app)
			if err != nil {
				return err
			}

			ctx := context.Background()
			group, ctx := errgroup.WithContext(ctx)
			for _, exe := range executions {
				if exe.Status != spin.Running {
					exeID := exe.ID
					group.Go(func() error {
						return orca.DeletePipeline(exeID)
					})
				}
			}

			err = group.Wait()
			return err
		},
	}

	return cleanCmd
}

func spinCmd() *cobra.Command {
	spinCmd := &cobra.Command{
		Use:     "spin",
		Aliases: []string{"spin"},
		Short:   "manage spin executions",
	}

	spinCmd.PersistentFlags().StringVarP(&orcaURL, "orca", "o", "", "orca url")
	spinCmd.PersistentFlags().StringVarP(&app, "app", "a", "", "application to clean")

	cleanCmd := clean()
	spinCmd.AddCommand(cleanCmd)

	return spinCmd
}
