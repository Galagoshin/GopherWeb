package framework

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Galagoshin/GoLogger/logger"
	"github.com/Galagoshin/GoUtils/configs"
	"github.com/Galagoshin/GoUtils/files"
	"github.com/Galagoshin/GoUtils/time"
	"github.com/Galagoshin/GopherWeb/web/plugins"
	"github.com/Galagoshin/GopherWeb/web/render"
	"github.com/Galagoshin/GopherWeb/web/server"
	"os"
	"os/exec"
	"strings"
)

type RuntimeMode uint8

const (
	ProductionMode = RuntimeMode(iota)
	DevelopMode    = RuntimeMode(iota)
)

var Mode RuntimeMode

var Config = &configs.Config{Name: "gopher"}
var BuildConfig = &configs.Config{Name: "build"}

func Build(err_out error) {
	build_str, exists := BuildConfig.Get("build")
	if !exists {
		logger.Panic(errors.New("\"run\" is not defined in the build config."))
	}
	cmd := exec.Command("go", strings.Split(build_str.(string), " ")...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err_out = err
		return
	}
	cmd.Stdin = os.Stdin
	err = cmd.Start()
	if err != nil {
		err_out = err
		return
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		logger.Print(m)
	}
	err = cmd.Wait()
	if err != nil {
		err_out = err
		return
	}
}

func IsUnderDocker() bool {
	file := files.File{Path: "/.dockerenv"}
	return file.Exists()
}

func HotReload() {
	logger.Print(fmt.Sprintf("HotReload finished (%f s.)", time.MeasureExecution(func() {
		plugins.DisableAllPlugins()
		Config = &configs.Config{Name: "gopher"}
		BuildConfig.Init(map[string]any{
			"run":   "run src/main.go",
			"build": "build -o gopher_server src/main.go",
		}, 1)
		Config.Init(map[string]any{
			"hot-reload-enabled":  "true",
			"live-reload-enabled": "true",
			"write-logs":          "true",
			"debug-level":         0,
		}, 1)
		server.Config = &configs.Config{Name: "server"}

		server.Stop()

		debugLevel, debugError := Config.Get("debug-level")
		if !debugError {
			logger.Panic(errors.New("\"debug-level\" is not defined in the framework config."))
		}

		writeLogs, writeLogsError := Config.Get("write-logs")
		if !writeLogsError {
			logger.Panic(errors.New("\"write-logs\" is not defined in the framework config."))
		}

		logger.SetDebugLevel(debugLevel.(int))
		logger.SetLogs(writeLogs.(string) == "true")

		render.LoadTemplates()

		plugins.EnableAllPlugins()

		server.Init()
		go server.Run()
	})))
}

func Shutdown(restart bool) {
	if restart {
		logger.Print("Framework is reloading...")
	} else {
		logger.Print("Framework is shuting down...")
	}
	plugins.DisableAllPlugins()
	err := server.Config.Save()
	if err != nil {
		logger.Error(err)
	}
	if restart {
		os.Exit(0)
	}
	os.Exit(130)
}
