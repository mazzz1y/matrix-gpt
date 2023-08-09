package main

import (
	"github.com/mazzz1y/matrix-gpt/internal/bot"
	"github.com/mazzz1y/matrix-gpt/internal/gpt"
	"github.com/urfave/cli/v2"
)

func run(c *cli.Context) error {
	mPassword := c.String("matrix-password")
	mUserId := c.String("matrix-id")
	mUrl := c.String("matrix-url")
	mRoom := c.String("matrix-room")
	sqlitePath := c.String("sqlite-path")

	gptModel := c.String("gpt-model")
	gptTimeout := c.Int("gpt-timeout")
	openaiToken := c.String("openai-token")
	maxAttempts := c.Int("max-attempts")

	historyExpire := c.Int("history-expire")
	historyLimit := c.Int("history-limit")
	userIDs := c.StringSlice("user-ids")

	logLevel := c.String("log-level")
	logType := c.String("log-type")

	setLogLevel(logLevel, logType)

	g := gpt.New(openaiToken, gptModel, historyLimit, gptTimeout, maxAttempts, userIDs)
	m, err := bot.NewBot(mUrl, mUserId, mPassword, sqlitePath, mRoom, historyExpire, g)
	if err != nil {
		return err
	}

	return m.StartHandler()
}
