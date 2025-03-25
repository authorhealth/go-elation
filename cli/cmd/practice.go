package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/authorhealth/go-elation"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var (
	listContactsNPI string
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

var getContactCmd = &cobra.Command{
	Use:  "get-contact [contact ID]",
	Args: cobra.ExactArgs(1),
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		contactID, _ := strconv.ParseInt(args[0], 10, 64)
		response, _, err := client.Contacts().Get(ctx, contactID)
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

var listContactsCmd = &cobra.Command{
	Use: "list-contacts",
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		response, _, err := client.Contacts().List(ctx, &elation.ListContactsOptions{
			Pagination: &elation.Pagination{
				Limit:  paginationLimit,
				Offset: paginationOffset,
			},
			NPI: listContactsNPI,
		})
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
	rootCmd.AddCommand(findPhysiciansCmd)
	rootCmd.AddCommand(getContactCmd)
	rootCmd.AddCommand(listContactsCmd)

	listContactsCmd.Flags().StringVar(&listContactsNPI, "npi", "", "")
}
