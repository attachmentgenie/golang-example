package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/attachmentgenie/golang-example/internal/server"
	promVersion "github.com/prometheus/common/version"
	"github.com/spf13/cobra"

	"github.com/attachmentgenie/golang-example/internal/server/api"
	"github.com/attachmentgenie/golang-example/internal/server/o11y"
)

var port int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start as a background process.",
	Long:  "Start as a background process.",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info(
			fmt.Sprintf("Starting up %s", Service),
			slog.String("version", promVersion.Version),
			slog.String("commit", promVersion.Revision),
		)

		apiCfg := server.Config{
			Port:    port,
			Service: Service,
		}
		mainSrv := api.NewServer(apiCfg)

		o11yCfg := server.Config{
			Port:    port+1,
			Service: Service,
		}
		o11ySrv := o11y.NewServer(o11yCfg)

		ctx, cancel := signal.NotifyContext(context.Background())
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT)

		go func() {
			slog.Info(
				"api server listening on...",
				slog.String("address", mainSrv.Addr),
			)
			if err := mainSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("error listening and serving api server", "error" , err.Error())
			}
		}()

		go func() {
			slog.Info(
				"o11y server listening on...",
				slog.String("address", o11ySrv.Addr),
			)
			if err := o11ySrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("error listening and serving o11y server", "error" , err.Error())
			}
		}()

		defer func() {
			if mainShutdownErr := mainSrv.Shutdown(ctx); mainShutdownErr != nil {
				slog.Error("error when shutting down the api server: ", "error", mainShutdownErr)
			}
			if mainO11yErr := o11ySrv.Shutdown(ctx); mainO11yErr != nil {
				slog.Error("error when shutting down the o11y server: ", "error", mainO11yErr)
			}
		}()

		sig := <-sigs

		fmt.Println(sig)

		cancel()

		slog.Info(fmt.Sprintf("%s has shutdown", Service))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&port, "port", 8088, fmt.Sprintf("port to expose %s on.", Service))
}
