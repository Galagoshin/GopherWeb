package web

import (
	"errors"
	"fmt"
	"github.com/Galagoshin/GoLogger/logger"
	"github.com/Galagoshin/GoUtils/files"
	"github.com/Galagoshin/GoUtils/time"
	"github.com/Galagoshin/GopherWeb/web/commands"
	"github.com/Galagoshin/GopherWeb/web/framework"
	"github.com/Galagoshin/GopherWeb/web/plugins"
	"github.com/Galagoshin/GopherWeb/web/render"
	"github.com/Galagoshin/GopherWeb/web/server"
	"github.com/Galagoshin/GopherWeb/web/tasks"
	"os"
	"strings"
)

const VERSION = "1.0.0-ALPHA1"

var commandsEnabled = false

func Init() {
	logger.Print(fmt.Sprintf("GopherWeb v%s has been loaded (%f s.)", VERSION, time.MeasureExecution(func() {

		getArgVal := func(arg string, args []string) (string, error) {
			for i := 0; i < len(args); i++ {
				if args[i] == arg {
					if i+1 >= len(args) {
						return "", errors.New(fmt.Sprintf("Arg \"%s\" has no value.", arg))
					}
					return args[i+1], nil
				}
			}
			return "", errors.New(fmt.Sprintf("Arg \"%s\" not found.", arg))
		}

		envFile := files.File{Path: ".env"}
		err := envFile.Open(os.O_RDWR)
		if err == nil {
			 for _, line := range strings.Split(envFile.ReadString(), "\n") {
				 keyval := strings.Split(line, "=")
				 key, val := keyval[0], keyval[1]
				 err = os.Setenv(key, val)
				 if err != nil {
					 logger.Panic(err)
				 }
			 }
		}

		modeArg, modeErr := getArgVal("--mode", os.Args)
		framework.Mode = framework.DevelopMode
		if framework.IsUnderDocker() {
			framework.Mode = framework.ProductionMode
		}
		if modeErr == nil {
			if modeArg == "prod" {
				framework.Mode = framework.ProductionMode
			} else if modeArg == "dev" {
				framework.Mode = framework.DevelopMode
			}
		}
		registerCommands := func() {
			logger.RegisterCommand(logger.Command{Name: "stop", Description: "Stop application.", Aliases: []string{"shutdown"}, Execute: commands.Stop})
			logger.RegisterCommand(logger.Command{Name: "help", Description: "Show all framework commands.", Aliases: []string{"commands"}, Execute: commands.Help})
			logger.RegisterCommand(logger.Command{Name: "reload", Description: "Do hot reload.", Execute: commands.Reload})
			logger.RegisterCommand(logger.Command{Name: "build", Description: "Build production version.", Aliases: []string{"compile"}, Execute: commands.Build})
			commandsEnabled = true
		}
		cmdsArg, cmdsErr := getArgVal("--commands", os.Args)
		if cmdsErr == nil {
			if cmdsArg == "enable" {
				registerCommands()
			}
		} else {
			if framework.Mode == framework.DevelopMode {
				registerCommands()
			}
		}

		framework.Config.Init(map[string]any{
			"hot-reload-enabled":  "true",
			"live-reload-enabled": "true",
			"write-logs":          "true",
			"debug-level":         0,
		}, 1)
		debugLevel, debugError := framework.Config.Get("debug-level")
		if !debugError {
			logger.Panic(errors.New("\"debug-level\" is not defined in the framework config."))
		}

		writeLogs, writeLogsError := framework.Config.Get("write-logs")
		if !writeLogsError {
			logger.Panic(errors.New("\"write-logs\" is not defined in the framework config."))
		}

		logger.SetDebugLevel(debugLevel.(int))
		logger.SetLogs(writeLogs.(string) == "true")

		framework.BuildConfig.Init(map[string]any{
			"run":   "run src/main.go",
			"build": "build -o gopher_server src/main.go",
		}, 1)

		hotreload, hotreloadError := framework.Config.Get("hot-reload-enabled")
		if !hotreloadError {
			logger.Panic(errors.New("\"hot-reload-enabled\" is not defined in the framework config."))
		}

		livereload, livereloadError := framework.Config.Get("live-reload-enabled")
		if !livereloadError {
			logger.Panic(errors.New("\"live-reload-enabled\" is not defined in the framework config."))
		}

		server.Init()
		logger.Print("Configuration has been loaded.")

		render.LoadTemplates()
		server.Static()

		if hotreload.(string) == "true" {
			tasks.InitHotReloadTask()
			tasks.HotReloadTask.Run(tasks.HotReloadTask)
		}

		plugins.EnableAllPlugins()

		if strings.Contains(VERSION, "ALPHA") {
			logger.Warning("This is a ALPHA version of GopherWeb framework, don't use it in production!")
		} else if strings.Contains(VERSION, "BETA") {
			logger.Warning("This is a BETA version of GopherWeb framework, don't use it in production!")
		}

		if livereload.(string) == "true" && framework.Mode == framework.DevelopMode {
			tasks.RestartTask.Run(tasks.RestartTask)
			logger.Warning("This app running in develop mode, don't use it in production!")
		}
	})))
	if commandsEnabled {
		logger.Info("Type \"help\" to get commands list.")
	}
}

func Run() {
	go server.Run()
	select {}
}
