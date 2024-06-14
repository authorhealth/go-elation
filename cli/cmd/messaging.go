package cmd

import (
	"context"
	"strconv"

	"github.com/authorhealth/go-elation"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var (
	findThreadMembersPatients []int64
	findThreadMembersUsers    []int64
)

var findThreadMembersCmd = &cobra.Command{
	Use: "find-thread-members",
	Run: wrapRunFunc(func(ctx context.Context, client *elation.Client, args []string) error {
		response, _, err := client.ThreadMemberSvc.Find(ctx, &elation.FindThreadMembersOptions{
			Patient: findThreadMembersPatients,
			User:    findThreadMembersUsers,
		})
		if err != nil {
			return err
		}

		spew.Dump(response)

		return nil
	}),
}

var getMessageThreadCmd = &cobra.Command{
	Use:  "get-message-thread [thread ID]",
	Args: cobra.ExactArgs(1),
	Run: wrapRunFunc(func(ctx context.Context, client *elation.Client, args []string) error {
		threadID, _ := strconv.ParseInt(args[0], 10, 64)
		response, _, err := client.MessageThreadSvc.Get(ctx, threadID)
		if err != nil {
			return err
		}

		spew.Dump(response)

		return nil
	}),
}

func init() {
	rootCmd.AddCommand(findThreadMembersCmd)
	rootCmd.AddCommand(getMessageThreadCmd)

	findThreadMembersCmd.Flags().Int64SliceVar(&findThreadMembersPatients, "patients", []int64{}, "")
	findThreadMembersCmd.Flags().Int64SliceVar(&findThreadMembersUsers, "users", []int64{}, "")
}
