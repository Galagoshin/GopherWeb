package tasks

import (
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

var lastDirHash, _ = crypto.HashDir("src", "HotReload", crypto.Hash1)

func RestartExecutor(args ...any) {
	task := args[0].(*scheduler.RepeatingTask)
	hash, err := crypto.HashDir("src", "HotReload", crypto.Hash1)
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
