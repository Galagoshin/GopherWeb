package tasks

import (
	"errors"
	"github.com/Galagoshin/GoLogger/logger"
	"github.com/Galagoshin/GoUtils/crypto"
	"github.com/Galagoshin/GoUtils/scheduler"
	"github.com/Galagoshin/GopherWeb/web/framework"
	time2 "time"
)

var RestartTask = &scheduler.RepeatingTask{
	Duration:   time2.Second,
	OnComplete: RestartExecutor,
}

var sourceDir string
var lastDirHash string

func RestartExecutor(args ...any) {
	task := args[0].(*scheduler.RepeatingTask)
	hash, err := crypto.HashDir(sourceDir, "HotReload", crypto.Hash1)
	if err != nil {
		logger.Error(err)
		task.Destroy()
	} else {
		if hash != lastDirHash {
			lastDirHash = hash
			framework.Shutdown(true)
		}
	}
}

func InitRestartTask() {
	source_key, exists := framework.BuildConfig.Get("source-dir")
	if !exists {
		logger.Panic(errors.New("\"source-dir\" is not defined in the framework config."))
	}
	sourceDir = source_key.(string)
	lastDirHash, _ = crypto.HashDir(sourceDir, "HotReload", crypto.Hash1)
}
