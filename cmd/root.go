package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	Service string
)

var rootCmd = &cobra.Command{
	Use:   Service,
	Short: "Basic golang http server example.",
	Long:  `Basic golang http server example.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
