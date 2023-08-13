package bot

import "strings"

const (
	emptyCommand         = ""
	generateImageCommand = "image"
	historyResetCommand  = "reset"
	helpCommand          = "help"
)

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

func commandIs(cmd, in string) bool {
	return cmd == in || string([]rune(cmd)[0]) == in
}
