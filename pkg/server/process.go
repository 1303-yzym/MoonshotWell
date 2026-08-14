package server

import (
	"fmt"
	"os"
	"os/user"

	"go.uber.org/zap"

	"github.com/fatih/color"
)

func PrintWorkingDirAndPID() {
	// 获取当前用户
	currentUser, err := user.Current()
	if err != nil {
		zap.L().Error("get the currently executing user", zap.Error(err))

		return
	}

	// 获取当前工作目录
	workingDir, err := os.Getwd()
	if err != nil {
		zap.L().Error("failed to get working directory:", zap.Error(err))

		return
	}

	// 获取当前进程ID
	pid := os.Getpid()
	// 格式化输出用户信息
	color.HiBlue(fmt.Sprintf("Username: %s GID: %s UID: %s\n", currentUser.Username, currentUser.Gid, currentUser.Uid))
	// 使用 color.HiBlue 输出带颜色的工作目录和PID信息
	color.HiBlue(fmt.Sprintf("WorkingDirectory: %s PID: %d", workingDir, pid))
}
