package o11y

import (
	"fmt"
	"net/http"

	"github.com/attachmentgenie/golang-example/internal/server"
	"github.com/prometheus/client_golang/prometheus"
	promCollectors "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promVersion "github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
)

func NewServer(cfg server.Config) http.Server {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		promCollectors.NewCollector(cfg.Service),
	)

	o11yMux := http.NewServeMux()
	o11yMux.Handle("/", O11yLandingPage(cfg.Service))
	o11yMux.Handle(
		"/metrics", promhttp.HandlerFor(
			reg,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			}),
	)
	return http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: http.NewCrossOriginProtection().Handler(o11yMux),
	}
}

func O11yLandingPage(service string) *web.LandingPageHandler {
	landingConfig := web.LandingConfig{
		Name: service,
		Links: []web.LandingLinks{
			{
				Address: "/metrics",
				Text:    "Metrics",
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
