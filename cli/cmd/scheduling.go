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

var getAppointment = &cobra.Command{
	Use:  "get-appointment [appointment ID]",
	Args: cobra.ExactArgs(1),
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		appointmentID, _ := strconv.ParseInt(args[0], 10, 64)
		response, _, err := client.Appointments().Get(ctx, appointmentID)
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

var updateAppointment = &cobra.Command{
	Use:  "update-appointment [appointment ID]",
	Args: cobra.ExactArgs(1),
	Run: wrapRunFunc(func(ctx context.Context, client elation.Client, args []string) error {
		appointmentID, _ := strconv.ParseInt(args[0], 10, 64)
		requestBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		request := &elation.AppointmentUpdate{}
		err = json.Unmarshal(requestBytes, request)
		if err != nil {
			return err
		}

		response, _, err := client.Appointments().Update(ctx, appointmentID, request)
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
	rootCmd.AddCommand(getAppointment)
	rootCmd.AddCommand(updateAppointment)
}
