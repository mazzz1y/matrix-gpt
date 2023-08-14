package bot

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"io"
	"net/http"
	"strings"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/attachment"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

const helpMessage = `**Commands**
- *!image [text]*: Creates an image based on the provided text.
- *!reset [text]*: Resets the user history. If a text is provided after the reset command, it will generate a GPT response based on this text.
- *[text]*: If only text is provided, the bot will generate a GPT-based response related to that text.

**Notes**
- You can use the first letter of a command as an alias. For example, "!i" for "!image".
- If you wish to terminate the current processing, simply delete your message from the chat.
- The bot responds with ❌ reaction if there are any errors. Contact the administrator if you see this.
`

type action func(context.Context, *user, *event.Event, string) error

// initBotActions is used to set up the possible actions the Bot can handle.
// This method should be called during the bot initialization process.
func (b *Bot) initBotActions() {
	b.actions = map[string]action{
		"":      b.completionResponse,
		"image": b.imageResponse,
		"reset": b.resetResponse,
		"help":  b.helpResponse,
	}
}

// getAction matches an input string to a bot action.
// It returns an exact match, the longest prefix match, or an unknownCommandError if no match is found.
func (b *Bot) getAction(input string) (action, error) {
	var bestMatch string
	for name := range b.actions {
		if input == name {
			return b.actions[input], nil
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
func (b *Bot) completionResponse(ctx context.Context, u *user, evt *event.Event, msg string) error {
	newHistory, err := b.gptClient.CreateCompletion(ctx, u.history.get(), msg)
	if err != nil {
		return err
	}

	u.history.save(newHistory)
	return b.markdownResponse(evt.RoomID, newHistory[len(newHistory)-1].Content)
}

// helpResponse responds with help message.
func (b *Bot) helpResponse(ctx context.Context, u *user, evt *event.Event, msg string) error {
	return b.markdownResponse(evt.RoomID, helpMessage)
}

// imageResponse responds to the user message with a DALL-E created image.
func (b *Bot) imageResponse(ctx context.Context, u *user, evt *event.Event, msg string) error {
	img, err := b.gptClient.CreateImage(ctx, msg)
	if err != nil {
		return err
	}

	imageBytes, err := getImageBytesFromURL(img.Data[0].URL)
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
func (b *Bot) markdownResponse(roomID id.RoomID, msg string) error {
	formattedMsg := format.RenderMarkdown(msg, true, false)
	_, err := b.client.SendMessageEvent(roomID, event.EventMessage, &formattedMsg)
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
