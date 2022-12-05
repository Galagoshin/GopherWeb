package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Galagoshin/GoLogger/logger"
	"github.com/Galagoshin/GoUtils/configs"
	"os"
	"os/exec"
	"strings"
)

func main() {
	logger.SetLogs(false)
	config := &configs.Config{Name: "build"}
	config.Init(map[string]any{
		"run":        "run src/main.go",
		"build":      "build -o gopher_server src/main.go",
		"source-dir": "src",
	}, 2)
	for {
		run_str, exists := config.Get("run")
		if !exists {
			logger.Panic(errors.New("\"run\" is not defined in the build config."))
		}
		cmd := exec.Command("go", strings.Split(run_str.(string), " ")...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			logger.Panic(err)
		}
		cmd.Stdin = os.Stdin
		var errbuf strings.Builder
		cmd.Stderr = &errbuf
		err = cmd.Start()
		if err != nil {
			logger.Panic(err)
		}
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
		}
		err = cmd.Wait()
		if err != nil {
			if err.Error() == "exit status 2" {
				stderr := errbuf.String()
				logger.Panic(errors.New("compilation error: \n" + stderr))
			} else {
				logger.Panic(err)
			}
			break
		}
	}
}
