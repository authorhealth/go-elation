package cmd

import (
	"context"
	"strconv"

	"github.com/authorhealth/go-elation"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var (
	findBillPatients     []int64
	findBillVisitNoteIDs []int64
)

var findBillsCmd = &cobra.Command{
	Use: "find-bills",
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		response, _, err := client.Bill().Find(ctx, &elation.FindBillOptions{
			Pagination: &elation.Pagination{
				Cursor: paginationCursor,
				Limit:  paginationLimit,
				Offset: paginationOffset,
			},
			Patient:     findBillPatients,
			VisitNoteID: findBillVisitNoteIDs,
		})
		if err != nil {
			return err
		}

		spew.Dump(response)

		return nil
	}),
}

var getBillCmd = &cobra.Command{
	Use:  "get-bill [bill ID]",
	Args: cobra.ExactArgs(1),
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		billID, _ := strconv.ParseInt(args[0], 10, 64)
		response, _, err := client.Bill().Get(ctx, billID)
		if err != nil {
			return err
		}

		spew.Dump(response)

		return nil
	}),
}

func init() {
	rootCmd.AddCommand(findBillsCmd)
	rootCmd.AddCommand(getBillCmd)

	findBillsCmd.Flags().Int64SliceVar(&findBillPatients, "patients", []int64{}, "")
	findBillsCmd.Flags().Int64SliceVar(&findBillVisitNoteIDs, "visit-note-ids", []int64{}, "")
}
