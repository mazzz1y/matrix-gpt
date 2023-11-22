package bot

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/h2non/filetype"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/attachment"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

type action func(context.Context, *user, *event.Event, string) error

// initBotActions is used to set up the possible actions the Bot can handle.
// This method should be called during the bot initialization process.
func (b *Bot) initBotActions() {
	b.actions = map[string]action{
		"":              b.completionResponse,
		"image-natural": b.imageResponse("natural"),
		"image-vivid":   b.imageResponse("vivid"),
		"reset":         b.resetResponse,
		"help":          b.helpResponse,
	}
}

// getAction matches an input string to a bot action.
// It returns an exact match, the longest prefix match, or an unknownCommandError if no match is found.
func (b *Bot) getAction(input string) (action, error) {
	var bestMatch string
	for name := range b.actions {
		if input == name || isAbbreviation(input, name) {
			return b.actions[name], nil
		}

		if strings.HasPrefix(name, input) && len(name) > len(bestMatch) {
			bestMatch = name
		}
	}

	if bestMatch == "" {
		return nil, &unknownCommandError{cmd: input}
	}

	return b.actions[bestMatch], nil
}

// completionResponse responds to a user message with a GPT-based completion.
// If the message is audio, it transcribes it before generating the response.
func (b *Bot) completionResponse(ctx context.Context, u *user, evt *event.Event, msg string) error {
	if evt.Content.AsMessage().MsgType == event.MsgAudio {
		fname, err := b.decryptAndStoreFile(evt)
		if err != nil {
			return err
		}
		defer os.Remove(fname)

		text, err := b.gptClient.CreateTranscription(ctx, fname)
		if err != nil {
			return err
		}

		msg = text
	}

	newHistory, err := b.gptClient.CreateCompletion(ctx, u.history.get(), msg)
	if err != nil {
		return err
	}

	u.history.save(newHistory)
	return b.markdownResponse(evt, false, newHistory[len(newHistory)-1].Content)
}

// helpResponse responds with help message.
func (b *Bot) helpResponse(ctx context.Context, u *user, evt *event.Event, msg string) error {
	return b.markdownResponse(evt, false, helpMsg)
}

// imageResponse responds to the user message with a DALL-E created image.
func (b *Bot) imageResponse(style string) action {
	return func(ctx context.Context, u *user, evt *event.Event, msg string) error {
		url, err := b.gptClient.CreateImage(ctx, style, msg)
		if err != nil {
			return err
		}

		imageBytes, err := getImageBytesFromURL(url)
		if err != nil {
			return err
		}

		cfg, err := png.DecodeConfig(bytes.NewReader(imageBytes))
		if err != nil {
			return err
		}

		content := b.createImageMessageContent(imageBytes, cfg)

		file := attachment.NewEncryptedFile()
		file.EncryptInPlace(imageBytes)

		req := mautrix.ReqUploadMedia{
			ContentBytes: imageBytes,
			ContentType:  "application/octet-stream",
		}

		upload, err := b.client.UploadMedia(req)
		if err != nil {
			return err
		}

		content.File = &event.EncryptedFileInfo{
			EncryptedFile: *file,
			URL:           upload.ContentURI.CUString(),
		}

		_, err = b.client.SendMessageEvent(evt.RoomID, event.EventMessage, content)
		return err
	}
}

// resetResponse clears the user's history. If a message is provided, it's processed as a new input.
// Otherwise, a reaction is sent to indicate successful history reset.
func (b *Bot) resetResponse(ctx context.Context, u *user, evt *event.Event, msg string) error {
	u.history.reset()
	if msg != "" {
		return b.completionResponse(ctx, u, evt, msg)
	} else {
		b.reactionResponse(evt, "✅")
	}
	return nil
}

// markdownResponse sends a message response in markdown format.
func (b *Bot) markdownResponse(evt *event.Event, reply bool, msg string) error {
	formattedMsg := format.RenderMarkdown(msg, true, false)
	if reply {
		formattedMsg.SetReply(evt)
	}

	_, err := b.client.SendMessageEvent(evt.RoomID, event.EventMessage, &formattedMsg)
	return err
}

// reactionResponse sends a reaction to a message.
func (b *Bot) reactionResponse(evt *event.Event, emoji string) {
	_, _ = b.client.SendReaction(evt.RoomID, evt.ID, emoji)
}

// markRead marks the given event as read by the bot.
func (b *Bot) markRead(evt *event.Event) {
	_ = b.client.MarkRead(evt.RoomID, evt.ID)
}

// startTyping notifies the room that the bot is typing.
func (b *Bot) startTyping(roomID id.RoomID) {
	_, _ = b.client.UserTyping(roomID, true, b.gptClient.GetTimeout())
}

// stopTyping notifies the room that the bot has stopped typing.
func (b *Bot) stopTyping(roomID id.RoomID) {
	_, _ = b.client.UserTyping(roomID, false, 0)
}

// createImageMessageContent creates the which contains the image information and the reply references.
func (b *Bot) createImageMessageContent(imageBytes []byte, cfg image.Config) *event.MessageEventContent {
	return &event.MessageEventContent{
		MsgType: event.MsgImage,
		Info: &event.FileInfo{
			Height:   cfg.Height,
			MimeType: http.DetectContentType(imageBytes),
			Width:    cfg.Height,
			Size:     len(imageBytes),
		},
	}
}

// getImageBytesFromURL returns the byte data from the image URL.
func getImageBytesFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// decryptAndStoreFile decrypts a file from an event and stores it locally in temp dir, returning the local file name.
func (b *Bot) decryptAndStoreFile(evt *event.Event) (string, error) {
	file := evt.Content.AsMessage().File

	if file == nil {
		return "", fmt.Errorf("no file found in message")
	}

	mxc, err := file.URL.Parse()
	if err != nil {
		return "", err
	}

	data, err := b.client.DownloadBytes(mxc)
	if err != nil {
		return "", err
	}

	err = file.DecryptInPlace(data)
	if err != nil {
		return "", err
	}

	return storeFile(data)
}

// storeFile stores a given byte slice as a file in a temporary directory and returns the generated file name.
func storeFile(data []byte) (fname string, err error) {
	ext, err := filetype.Match(data)
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp("", "*."+ext.Extension)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	_, err = tempFile.Write(data)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

// isAbbreviation checks if the input is a valid abbreviation of the action name
// by matching the input with the initials of hyphen-separated parts in the action name.
func isAbbreviation(input, actionName string) bool {
	nameParts := strings.Split(actionName, "-")

	if len(input) != len(nameParts) {
		return false
	}

	for i, part := range nameParts {
		if !strings.HasPrefix(part, input[i:i+1]) {
			return false
		}
	}

	return true
}
