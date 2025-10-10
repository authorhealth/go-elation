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

var getInsurancePoliciesActiveOnly bool

var getInsurancePoliciesCmd = &cobra.Command{
	Use:  "get-insurance-policies [patient ID]",
	Args: cobra.ExactArgs(1),
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		patientID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("parsing patient ID: %w", err)
		}

		response, _, err := client.InsurancePolicies().Find(ctx, patientID, &elation.FindInsurancePoliciesOptions{})
		if err != nil {
			return err
		}

		spew.Dump(response)

		return nil
	}),
}

var getInsurancePolicyCmd = &cobra.Command{
	Use:  "get-insurance-policy [patient ID] [policy ID]",
	Args: cobra.ExactArgs(2),
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		patientID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("parsing patient ID: %w", err)
		}

		policyID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("parsing policy ID: %w", err)
		}

		response, _, err := client.InsurancePolicies().Get(ctx, patientID, policyID)
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
	getInsurancePoliciesCmd.Flags().BoolVar(&getInsurancePoliciesActiveOnly, "active-only", false, "Include active policies only")

	rootCmd.AddCommand(getInsurancePoliciesCmd)
	rootCmd.AddCommand(getInsurancePolicyCmd)
}
