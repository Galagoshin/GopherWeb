package commands

import (
	"errors"
	"github.com/Galagoshin/GoLogger/logger"
	"github.com/Galagoshin/GopherWeb/web/framework"
	"github.com/Galagoshin/GopherWeb/web/tasks"
)

func Reload(string, []string) {
	prevLive, exists := framework.Config.Get("live-reload-enabled")
	if !exists {
		logger.Panic(errors.New("\"live-reload-enabled\" was removed from the framework config."))
	}
	prevHot, exists := framework.Config.Get("hot-reload-enabled")
	if !exists {
		logger.Panic(errors.New("\"hot-reload-enabled\" was removed from the framework config."))
	}
	framework.HotReload()
	hotreload, hotreloadError := framework.Config.Get("hot-reload-enabled")
	if !hotreloadError {
		logger.Panic(errors.New("\"hot-reload-enabled\" is not defined in the framework config."))
	}
	if hotreload.(string) == "false" && prevHot.(string) != hotreload.(string) {
		tasks.HotReloadTask.Destroy()
	} else if hotreload.(string) == "true" && prevHot.(string) != hotreload.(string) {
		tasks.InitHotReloadTask()
		tasks.HotReloadTask.Run(tasks.HotReloadTask)
	}
	livereload, livereloadError := framework.Config.Get("live-reload-enabled")
	if !livereloadError {
		logger.Panic(errors.New("\"live-reload-enabled\" is not defined in the framework config."))
	}
	if livereload.(string) == "false" && prevLive.(string) != livereload.(string) {
		tasks.RestartTask.Destroy()
	} else if livereload.(string) == "true" && prevLive.(string) != livereload.(string) {
		tasks.RestartTask.Run(tasks.RestartTask)
	}
}
