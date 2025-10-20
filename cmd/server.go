package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/prometheus/client_golang/prometheus"
	promCollectors "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promVersion "github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/spf13/cobra"
)

var port int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start as a background process.",
	Long:  "Start as a background process.",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info(
			"Starting up...",
			slog.String("version", promVersion.Version),
			slog.String("commit", promVersion.Revision),
		)
		mcpServer := mcp.NewServer(&mcp.Implementation{
			Name:    Service,
			Version: promVersion.Version,
		}, nil)

		reg := prometheus.NewRegistry()
		reg.MustRegister(
			promCollectors.NewCollector(Service),
		)

		http.Handle("/", landingPage())
		http.HandleFunc("/health", pingHTTP)
		mcp.AddTool(mcpServer, &mcp.Tool{Name: "health"}, pingMcp)
		http.Handle(
			"/mcp", mcp.NewStreamableHTTPHandler(
				func(req *http.Request) *mcp.Server {
					return mcpServer
				},
				nil),
		)
		http.Handle(
			"/metrics", promhttp.HandlerFor(
				reg,
				promhttp.HandlerOpts{
					EnableOpenMetrics: true,
				}),
		)
		http.HandleFunc("/ready", pingHTTP)
		err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
		if err != nil {
			slog.Error(err.Error(), err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&port, "port", 8088, "port to expose service on.")
}

func ping() string {
	return "pong"
}
func pingHTTP(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(w, ping())
	if err != nil {
		return
	}
}

func pingMcp(ctx context.Context, req *mcp.CallToolRequest, input any) (
	*mcp.CallToolResult,
	any,
	error,
) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: ping()},
		},
	}, nil, nil
}

func landingPage() *web.LandingPageHandler {
	landingConfig := web.LandingConfig{
		Name:    Service,
		Version: promVersion.Version,
		Links: []web.LandingLinks{
			{
				Address: "/health",
				Text:    "Health",
			},
			{
				Address: "/mcp",
				Text:    "MCP server",
			},
			{
				Address: "/metrics",
				Text:    "Metrics",
			},
			{
				Address: "/ready",
				Text:    "Ready",
			},
		},
	}
	landingPage, err := web.NewLandingPage(landingConfig)
	if err != nil {
		panic(err)
	}

	return landingPage
}
