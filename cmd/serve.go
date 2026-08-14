package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",
	Run: func(cmd *cobra.Command, args []string) {
		// We get the configuration value from Viper, not from the flag directly.
		port := Cfg.Server.HTTP.Port
		fmt.Printf("Starting server on port: %d \n", port)
		// In a real app, you would start a server here.
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
