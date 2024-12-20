package cmd

import (
	"context"
	"strconv"

	"github.com/authorhealth/go-elation"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var getDiscontinuedMedication = &cobra.Command{
	Use:  "get-discontinued-medication [discontinued medication ID]",
	Args: cobra.ExactArgs(1),
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		discontinuedMedicationID, _ := strconv.ParseInt(args[0], 10, 64)
		response, _, err := client.DiscontinuedMedications().Get(ctx, discontinuedMedicationID)
		if err != nil {
			return err
		}

		spew.Dump(response)

		return nil
	}),
}

var findDiscontinuedMedications = &cobra.Command{
	Use: "find-discontinued-medications",
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		response, _, err := client.DiscontinuedMedications().Find(ctx, &elation.FindDiscontinuedMedicationsOptions{
			Pagination: &elation.Pagination{
				Limit:  paginationLimit,
				Offset: paginationOffset,
			},
		})
		if err != nil {
			return err
		}

		spew.Dump(response)

		return nil
	}),
}

var getMedication = &cobra.Command{
	Use:  "get-medication [medication ID]",
	Args: cobra.ExactArgs(1),
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		medicationID, _ := strconv.ParseInt(args[0], 10, 64)
		response, _, err := client.Medications().Get(ctx, medicationID)
		if err != nil {
			return err
		}

		spew.Dump(response)

		return nil
	}),
}

var findPrescriptionFills = &cobra.Command{
	Use: "find-prescription-fills",
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		response, _, err := client.PrescriptionFills().Find(ctx, &elation.FindPrescriptionFillsOptions{
			Pagination: &elation.Pagination{
				Limit:  paginationLimit,
				Offset: paginationOffset,
			},
		})
		if err != nil {
			return err
		}

		spew.Dump(response)

		return nil
	}),
}

var getPrescriptionFill = &cobra.Command{
	Use:  "get-prescription-fill [prescription fill ID]",
	Args: cobra.ExactArgs(1),
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		fillID, _ := strconv.ParseInt(args[0], 10, 64)
		response, _, err := client.PrescriptionFills().Get(ctx, fillID)
		if err != nil {
			return err
		}

		spew.Dump(response)

		return nil
	}),
}

func init() {
	rootCmd.AddCommand(findDiscontinuedMedications)
	rootCmd.AddCommand(findPrescriptionFills)
	rootCmd.AddCommand(getDiscontinuedMedication)
	rootCmd.AddCommand(getMedication)
	rootCmd.AddCommand(getPrescriptionFill)
}
