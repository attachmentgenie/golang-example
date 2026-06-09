package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/attachmentgenie/golang-example/internal/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	promVersion "github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
)

func NewServer(cfg server.Config) http.Server {
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    cfg.Service,
		Version: promVersion.Version,
	}, nil)

	mainMux := http.NewServeMux()
	mainMux.Handle("/", MainLandingPage(cfg.Service))
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
	return http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: http.NewCrossOriginProtection().Handler(mainMux),
	}
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

func MainLandingPage(service string) *web.LandingPageHandler {
	landingConfig := web.LandingConfig{
		Name: service,
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
