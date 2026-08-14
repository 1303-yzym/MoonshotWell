package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/1303-yzym/MoonshotWell/internal/infrastructure/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	cfgFilePath string
	logDirPath  string
	env         string
	ShouldRun   bool
	Cfg         *config.Config

	rootCmd = &cobra.Command{
		Use:   SERVERNAME,
		Short: DESCRIPTION,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := config.InitConfig(cfgFilePath); err != nil {
				log.Fatalf("Failed to initialize configuration: %v", err)
			}

			Cfg = config.Load()
			ShouldRun = true
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.L().Error("rootCmd.Execute error", zap.Error(err))
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SetVersionTemplate(fmt.Sprintf("version: %s revision: %s", VERSION, REVISION))
	rootCmd.Flags().BoolP("help", "h", false, "get help")
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")

	// 默认配置文件存放于工作目录下
	rootCmd.PersistentFlags().StringVarP(&cfgFilePath, "config", "c", "./config.yaml", "Set the path to the configuration file")
	rootCmd.Flags().StringVarP(&logDirPath, "log", "l", "./logs", "Set the path of the journal file")
	rootCmd.Flags().StringVarP(&env, "env", "e", "", "Overwrite configuration file env")
	_ = viper.BindPFlag("log.log_dir", rootCmd.Flags().Lookup("log"))
	_ = viper.BindPFlag("env", rootCmd.Flags().Lookup("env"))
}
