package main

import (
	"fmt"
	"os"

	"github.com/1303-yzym/MoonshotWell/internal/infrastructure/config"
	"github.com/1303-yzym/MoonshotWell/internal/infrastructure/state"
	"github.com/1303-yzym/MoonshotWell/pkg/infra"
	"github.com/1303-yzym/MoonshotWell/pkg/server"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	cfgFilePath string
	logDirPath  string
	env         string
	ShouldRun   bool
	appState    *state.AppState

	rootCmd = &cobra.Command{
		Use:   SERVERNAME,
		Short: DESCRIPTION,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// 初始化配置
			if err := config.InitConfig(cfgFilePath); err != nil {
				zap.L().Fatal("Failed to initialize configuration", zap.Error(err))
			}

			cfg := config.Load()
			//
			infra.InitLogger(cfg.Log, cfg.IsDev(),
				infra.ServerInfo{
					ServerName:  SERVERNAME,
					ServiceName: cfg.ServiceName,
					Version:     VERSION,
					ReVersion:   REVISION,
				}.LogField()...,
			)

			// 初始化全局句柄
			appState = state.InitAppState()
			// 允许启动
			ShouldRun = true
		},
	}
)

func execute() {
	color.HiBlue(fmt.Sprintf("Version: %s ReVision: %s\n", VERSION, REVISION))

	server.PrintWorkingDirAndPID()

	if err := rootCmd.Execute(); err != nil {
		zap.L().Error("rootCmd.Execute error", zap.Error(err))
		os.Exit(1)
	}
}

func init() {
	cobra.EnableTraverseRunHooks = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SetVersionTemplate(fmt.Sprintf("version: %s revision: %s", VERSION, REVISION))
	rootCmd.Flags().BoolP("help", "h", false, "get help")
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")

	// 默认配置文件存放于工作目录下
	rootCmd.PersistentFlags().StringVarP(&cfgFilePath, "config", "c", "./config.yaml", "Set the path to the configuration file")
	rootCmd.Flags().StringVarP(&logDirPath, "log", "l", "./logs", "Set the path of the journal file")
	// PersistentFlags()为全局变量能够延申到子命令使用
	rootCmd.Flags().StringVarP(&env, "env", "e", "", "Overwrite configuration file env")
	_ = viper.BindPFlag("log.log_dir", rootCmd.Flags().Lookup("log"))
	_ = viper.BindPFlag("env", rootCmd.Flags().Lookup("env"))
}
