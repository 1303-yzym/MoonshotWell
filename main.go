/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"os"

	"github.com/1303-yzym/MoonshotWell/cmd"
)

func main() {
	cmd.Execute()
	if !cmd.ShouldRun {
		os.Exit(0)
	}

	log.Printf("配置文件：%+v", cmd.Cfg)
}
