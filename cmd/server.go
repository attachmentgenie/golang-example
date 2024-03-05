package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	promversion "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/spf13/cobra"
)

var port int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start as a background process.",
	Long:  "Start as a background process.",
	Run: func(cmd *cobra.Command, args []string) {
		reg := prometheus.NewRegistry()
		reg.MustRegister(
			promversion.NewCollector(Service),
		)

		http.Handle("/", landingPage())
		http.HandleFunc("/health", ping)
		http.Handle(
			"/metrics", promhttp.HandlerFor(
				reg,
				promhttp.HandlerOpts{
					EnableOpenMetrics: true,
				}),
		)
		http.HandleFunc("/ready", ping)
		log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&port, "port", 8088, "port to expose service on.")
}

func ping(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(w, "pong")
	if err != nil {
		return
	}
}

func landingPage() *web.LandingPageHandler {
	landingConfig := web.LandingConfig{
		Name:    Service,
		Version: version.Version,
		Links: []web.LandingLinks{
			{
				Address: "/health",
				Text:    "Health",
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
