package main

import (
	"context"
	"os"

	"github.com/1303-yzym/MoonshotWell/cmd"
	"github.com/1303-yzym/MoonshotWell/pkg/infra/logger"
	"github.com/1303-yzym/MoonshotWell/pkg/signal"
	"go.uber.org/zap"
)

func main() {
	cmd.Execute()
	if !cmd.ShouldRun {
		os.Exit(0)
	}

	// 程序退出
	signal.On(func(ctx context.Context) {
		log := logger.App()
		log.Info("stop server...")

		if err := cmd.Adapter.Shutdown(ctx); err != nil {
			log.Error("shutdown err", zap.Error(err))
		}

		log.Info("stop server done")

	})

	signal.Listen()
}
