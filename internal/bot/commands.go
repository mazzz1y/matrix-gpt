package bot

import (
	"fmt"
	"strings"
)

type unknownCommandError struct {
	cmd string
}

func (e *unknownCommandError) Error() string {
	return fmt.Sprintf("command '!%s' does not exist", e.cmd)
}

func extractCommand(s string) (cmd string) {
	if strings.HasPrefix(s, "!") && len(s) > 1 {
		//Get the word after '!'
		command := strings.Fields(s)[0][1:]
		return command
	}
	return ""
}

func trimCommand(s string) string {
	if strings.HasPrefix(s, "!") && len(s) > 1 {
		//Remove command from s and clean up leading spaces
		trimmed := strings.TrimPrefix(s, strings.Fields(s)[0])
		return strings.TrimSpace(trimmed)
	}
	return s
}
