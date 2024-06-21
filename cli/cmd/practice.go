package cmd

import (
	"context"

	"github.com/authorhealth/go-elation"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var findPhysiciansCmd = &cobra.Command{
	Use: "find-physicians",
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		response, _, err := client.Physicians().Find(ctx, &elation.FindPhysiciansOptions{})
		if err != nil {
			return err
		}

		spew.Dump(response)

		return nil
	}),
}

func init() {
	rootCmd.AddCommand(findPhysiciansCmd)
}
