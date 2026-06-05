package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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
		slog.Info(
			"Listening on...",
			slog.String("port", strconv.Itoa(port)),
			slog.String("o11y", strconv.Itoa(port+1)),
		)
		mcpServer := mcp.NewServer(&mcp.Implementation{
			Name:    Service,
			Version: promVersion.Version,
		}, nil)

		reg := prometheus.NewRegistry()
		reg.MustRegister(
			promCollectors.NewCollector(Service),
		)

		mainMux := http.NewServeMux()
		mainMux.Handle("/", landingPage())
		mainMux.HandleFunc("/health", pingHTTP)
		mcp.AddTool(mcpServer, &mcp.Tool{Name: "health"}, pingMcp)
		mainMux.Handle(
			"/mcp", mcp.NewStreamableHTTPHandler(
				func(req *http.Request) *mcp.Server {
					return mcpServer
				},
				nil),
		)
		mainMux.HandleFunc("/ready", pingHTTP)
		mainSrv := &http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: http.NewCrossOriginProtection().Handler(mainMux),
		}

		o11yMux := http.NewServeMux()
		o11yMux.Handle(
			"/metrics", promhttp.HandlerFor(
				reg,
				promhttp.HandlerOpts{
					EnableOpenMetrics: true,
				}),
		)
		o11ySrv := &http.Server{
			Addr:    fmt.Sprintf(":%v", port+1),
			Handler: http.NewCrossOriginProtection().Handler(o11yMux),
		}

		ctx, cancel := context.WithCancel(context.Background())
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT)

		go func() {
			err := mainSrv.ListenAndServe()
			if err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
		}()

		go func() {
			err := o11ySrv.ListenAndServe()
			if err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
		}()

		defer func() {
			if err := mainSrv.Shutdown(ctx); err != nil {
				slog.Error("error when shutting down the main server: ", "error", err)
			}
			if err := o11ySrv.Shutdown(ctx); err != nil {
				slog.Error("error when shutting down the o11y server: ", "error", err)
			}
		}()

		sig := <-sigs
		fmt.Println(sig)

		cancel()

		slog.Error("service has shutdown")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&port, "port", 8088, "port to expose service on.")
}

func ping(protocol string) string {
	return fmt.Sprintf("pong %s", protocol)
}
func pingHTTP(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(w, "%s", ping("HTTP"))
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
			&mcp.TextContent{Text: ping("MCP")},
		},
	}, nil, nil
}

func landingPage() *web.LandingPageHandler {
	landingConfig := web.LandingConfig{
		Name: Service,
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
		Profiling: "false",
		Version:   promVersion.Version,
	}
	landingPage, err := web.NewLandingPage(landingConfig)
	if err != nil {
		panic(err)
	}

	return landingPage
}
