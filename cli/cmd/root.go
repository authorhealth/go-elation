package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/authorhealth/go-elation"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "elation",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type runFunc func(ctx context.Context, client *elation.Client, args []string) error

func wrapRunFunc(runFunc runFunc) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load()

		ctx := context.Background()

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)

		client := elation.NewClient(&http.Client{},
			os.Getenv("ELATION_TOKEN_URL"),
			os.Getenv("ELATION_CLIENT_ID"),
			os.Getenv("ELATION_CLIENT_SECRET"),
			os.Getenv("ELATION_BASE_URL"))

		err := runFunc(ctx, client, args)
		if err != nil {
			slog.ErrorContext(ctx, "error running command", slog.Any("error", err))
		}
	}
}
