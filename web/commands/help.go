package commands

import (
	"fmt"
	"github.com/Galagoshin/GoLogger/logger"
)

func Help(command string, args []string) {
	help_info := "List of commands:\n"
	for alias, cmd := range logger.GetAllCommands() {
		help_info += fmt.Sprintf("%s - %s\n", alias, cmd.Description)
	}
	logger.CommandInfo(help_info[:len(help_info)-1])
}
