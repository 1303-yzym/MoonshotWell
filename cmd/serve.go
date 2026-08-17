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
			// TODO 初始化应用层
			app = application.InitApplication(appState)
			// TODO 加载出口适配器
			Adapter = adapter.LoadAdapter(appState, app)
		},
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)
}
