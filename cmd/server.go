package cmd

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promversion "github.com/prometheus/common/version"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

var port int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start as a background process.",
	Long:  "Start as a background process.",
	Run: func(cmd *cobra.Command, args []string) {
		reg := prometheus.NewRegistry()
		reg.MustRegister(
			promversion.NewCollector(service),
		)

		http.HandleFunc("/", ping)
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
	fmt.Fprintf(w, "pong")
}
