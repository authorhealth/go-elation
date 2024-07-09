package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/authorhealth/go-elation"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	paginationLimit  int
	paginationOffset int
)

var rootCmd = &cobra.Command{
	Use: "elation",
}

func init() {
	rootCmd.PersistentFlags().IntVar(&paginationLimit, "limit", 0, "")
	rootCmd.PersistentFlags().IntVar(&paginationOffset, "offset", 0, "")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type runFunc func(ctx context.Context, client elation.Client, args []string) error

func wrapRunFunc(runFunc runFunc) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load()

		ctx := context.Background()

		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		slog.SetDefault(logger)

		client := elation.NewHTTPClient(
			&http.Client{
				Timeout: 15 * time.Second,
			},
			os.Getenv("ELATION_TOKEN_URL"),
			os.Getenv("ELATION_CLIENT_ID"),
			os.Getenv("ELATION_CLIENT_SECRET"),
			os.Getenv("ELATION_BASE_URL"))

		err := runFunc(ctx, client, args)
		if err != nil {
			apiError := &elation.Error{}
			if errors.As(err, &apiError) {
				slog.ErrorContext(ctx, "API error running command",
					slog.Any("error", err),
					slog.Int("statusCode", apiError.StatusCode),
					slog.String("body", apiError.Body))
			} else {
				slog.ErrorContext(ctx, "error running command", slog.Any("error", err))
			}
		}
	}
}
