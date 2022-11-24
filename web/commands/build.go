package commands

import (
	"fmt"
	"github.com/Galagoshin/GoLogger/logger"
	"github.com/Galagoshin/GoUtils/time"
	"github.com/Galagoshin/GopherWeb/web/framework"
)

func Build(string, []string) {
	var err error
	compile_time := time.MeasureExecution(func() {
		framework.Build(err)
	})
	if err == nil {
		logger.Print(fmt.Sprintf("Compilation done (%f s.)", compile_time))
	} else {
		logger.Error(err)
	}
}
