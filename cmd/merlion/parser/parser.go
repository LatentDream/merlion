// Package parser contains the logic to parse the command line arguments
package parser

import (
	"strings"

	"github.com/charmbracelet/log"
)

func GetArg(args []string, helpFunc func(bool)) (string, []string) {
	if len(args) == 0 {
		helpFunc(true)
	}
	currArg := args[0]
	args = args[1:]
	log.Debugf("currArg: %s", currArg)
	return currArg, args
}

// ParseArgs separates flags from commands and returns both
func SplitCmdsAndFlags(args []string) (flags []string, commands []string) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			flags = append(flags, arg)
		} else {
			commands = append(commands, arg)
		}
	}
	return flags, commands
}
