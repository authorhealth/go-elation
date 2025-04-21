package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/authorhealth/go-elation"
	"github.com/spf13/cobra"
)

var getPharmacyCmd = &cobra.Command{
	Use:  "get-pharmacy [pharmacy NCPDPID]",
	Args: cobra.ExactArgs(1),
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		ncpdpid := args[0]
		response, _, err := client.Pharmacies().Get(ctx, ncpdpid)
		if err != nil {
			return err
		}

		responseJson, err := json.Marshal(response)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(os.Stdout, string(responseJson))

		return err
	}),
}

func init() {
	rootCmd.AddCommand(getPharmacyCmd)
}
