package main

import (
	adp "github.com/1303-yzym/MoonshotWell/internal/adapter"
	"github.com/1303-yzym/MoonshotWell/internal/application"
	"github.com/spf13/cobra"
)

var (
	app     *application.Application
	adapter *adp.Adapter

	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Starts the server",
		Run: func(cmd *cobra.Command, args []string) {
			// 初始化应用层
			app = application.InitApplication(appState)
			// 加载出口适配器
			adapter = adp.LoadAdapter(appState, app)
		},
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)
}
