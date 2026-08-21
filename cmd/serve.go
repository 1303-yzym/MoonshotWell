package cmd

import (
	"github.com/1303-yzym/MoonshotWell/internal/adapter"
	"github.com/1303-yzym/MoonshotWell/internal/application"
	"github.com/spf13/cobra"
)

var (
	app     *application.Application
	Adapter *adapter.Adapter

	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Starts the server",
		Run: func(cmd *cobra.Command, args []string) {
			app = application.InitApplication(appState)
			Adapter = adapter.LoadAdapter(appState, app)
		},
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)
}
