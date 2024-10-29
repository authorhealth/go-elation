package cmd

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/authorhealth/go-elation"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var deleteSubscription = &cobra.Command{
	Use:  "delete-subscription [subscription ID]",
	Args: cobra.ExactArgs(1),
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		subscriptionID, _ := strconv.ParseInt(args[0], 10, 64)
		_, err := client.Subscriptions().Delete(ctx, subscriptionID)
		if err != nil {
			return err
		}

		return nil
	}),
}

var findSubscriptions = &cobra.Command{
	Use: "find-subscriptions",
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		response, _, err := client.Subscriptions().Find(ctx)
		if err != nil {
			return err
		}

		spew.Dump(response)

		return nil
	}),
}

var subscribe = &cobra.Command{
	Use: "subscribe",
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		requestBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		request := &elation.Subscribe{}
		err = json.Unmarshal(requestBytes, request)
		if err != nil {
			return err
		}

		response, _, err := client.Subscriptions().Subscribe(ctx, request)
		if err != nil {
			return err
		}

		spew.Dump(response)

		return nil
	}),
}

func init() {
	rootCmd.AddCommand(deleteSubscription)
	rootCmd.AddCommand(findSubscriptions)
	rootCmd.AddCommand(subscribe)
}
