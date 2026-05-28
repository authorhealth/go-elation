package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/authorhealth/go-elation"
	"github.com/mark3labs/mcp-go/server"
)

const (
	defaultHTTPTimeout = 15 * time.Second
	serverName         = "go-elation-mcp"
	serverVersion      = "0.2.0"
)

type serverOptions struct {
	allowUnsafeTools bool
}

func main() {
	opts, err := parseServerOptions(os.Args[1:])
	if err != nil {
		slog.Error("failed to parse server options", "error", err)
		os.Exit(2)
	}

	client := elation.NewHTTPClient(
		&http.Client{Timeout: defaultHTTPTimeout},
		os.Getenv("ELATION_TOKEN_URL"),
		os.Getenv("ELATION_CLIENT_ID"),
		os.Getenv("ELATION_CLIENT_SECRET"),
		os.Getenv("ELATION_BASE_URL"),
	)

	s := server.NewMCPServer(serverName, serverVersion)
	registerTools(s, client, opts.allowUnsafeTools)

	if err := server.ServeStdio(s); err != nil {
		if errors.Is(err, context.Canceled) {
			slog.Info("server stopped")
			os.Exit(0)
		}

		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func parseServerOptions(args []string) (serverOptions, error) {
	var opts serverOptions

	fs := flag.NewFlagSet(serverName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.BoolVar(&opts.allowUnsafeTools, "allow-unsafe-tools", false, "Expose mutative tools (create/update/delete)")

	if err := fs.Parse(args); err != nil {
		return serverOptions{}, err
	}

	return opts, nil
}
