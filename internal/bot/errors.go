package bot

import "fmt"

type unknownCommandError struct {
	cmd string
}

func (e *unknownCommandError) Error() string {
	return fmt.Sprintf("command '!%s' does not exist", e.cmd)
}
