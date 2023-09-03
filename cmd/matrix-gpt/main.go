package main

import (
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
)

var version = "git"

func main() {
	app := &cli.App{
		Name:    "matrix-gpt",
		Version: version,
		Usage:   "GPT Matrix Bot",
		Action:  run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "matrix-password",
				Usage:    "Matrix password",
				EnvVars:  []string{"MATRIX_PASSWORD"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "matrix-id",
				Usage:    "Matrix user ID",
				EnvVars:  []string{"MATRIX_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "matrix-url",
				Usage:    "Matrix server URL",
				EnvVars:  []string{"MATRIX_URL"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "openai-token",
				Usage:    "OpenAI API token",
				EnvVars:  []string{"OPENAI_TOKEN"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "sqlite-path",
				Usage:    "Path to SQLite database",
				EnvVars:  []string{"SQLITE_PATH"},
				Required: true,
			},
			&cli.IntFlag{
				Name:    "history-limit",
				Usage:   "Maximum number of history entries",
				EnvVars: []string{"HISTORY_LIMIT"},
				Value:   5,
			},
			&cli.IntFlag{
				Name:    "history-expire",
				Usage:   "Time after which history entries expire (in hours)",
				EnvVars: []string{"HISTORY_EXPIRE"},
				Value:   3,
			},
			&cli.StringFlag{
				Name:    "gpt-model",
				Usage:   "GPT model name/version",
				EnvVars: []string{"GPT_MODEL"},
				Value:   openai.GPT3Dot5Turbo,
			},
			&cli.IntFlag{
				Name:    "gpt-timeout",
				Usage:   "Time to wait for a GPT response (in seconds)",
				EnvVars: []string{"GPT_TIMEOUT"},
				Value:   120,
			},
			&cli.IntFlag{
				Name:    "max-attempts",
				Usage:   "Maximum number of attempts for GPT requests",
				EnvVars: []string{"MAX_ATTEMPTS"},
				Value:   3,
			},
			&cli.StringSliceFlag{
				Name:     "user-ids",
				Usage:    "List of allowed Matrix user IDs",
				EnvVars:  []string{"USER_IDS"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "log-level",
				Value:   "info",
				Usage:   "Logging level (e.g. debug, info, warn, error, fatal, panic, no)",
				EnvVars: []string{"LOG_LEVEL"},
			},
			&cli.StringFlag{
				Name:    "log-type",
				Value:   "pretty",
				Usage:   "Logging format/type (e.g. pretty, json)",
				EnvVars: []string{"LOG_TYPE"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("\n" + err.Error())
	}
}
