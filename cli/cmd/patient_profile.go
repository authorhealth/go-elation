package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/authorhealth/go-elation"
	"github.com/spf13/cobra"
)

var getPatient = &cobra.Command{
	Use:  "get-patient [patient ID]",
	Args: cobra.ExactArgs(1),
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		patientID, _ := strconv.ParseInt(args[0], 10, 64)
		response, _, err := client.Patients().Get(ctx, patientID)
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

var updatePatient = &cobra.Command{
	Use:  "update-patient [patient ID]",
	Args: cobra.ExactArgs(1),
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		patientID, _ := strconv.ParseInt(args[0], 10, 64)
		requestBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		request := &elation.PatientUpdate{}
		err = json.Unmarshal(requestBytes, request)
		if err != nil {
			return err
		}

		response, _, err := client.Patients().Update(ctx, patientID, request)
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
	rootCmd.AddCommand(updatePatient)
	rootCmd.AddCommand(getPatient)
}
