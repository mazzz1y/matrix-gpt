package bot

import (
	"bytes"
	"github.com/mazzz1y/matrix-gpt/internal/gpt"
	"image"
	"image/png"
	"io"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/attachment"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
	"net/http"
)

// completionResponse responds to a user message with a GPT-based completion.
func (b *Bot) completionResponse(u *gpt.User, roomID id.RoomID, msg string) error {
	answer, err := b.gptClient.CreateCompletion(u, msg)
	if err != nil {
		return err
	}

	return b.markdownResponse(roomID, answer)
}

// helpResponse responds with help message.
func (b *Bot) helpResponse(roomID id.RoomID) error {
	return b.markdownResponse(roomID, helpMessage)
}

// imageResponse responds to the user message with a DALL-E created image.
func (b *Bot) imageResponse(roomID id.RoomID, msg string) error {
	img, err := b.gptClient.CreateImage(msg)
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

	_, err = b.client.SendMessageEvent(roomID, event.EventMessage, content)
	return err
}

// resetResponse clears the user's history. If a message is provided, it's processed as a new input.
// Otherwise, a reaction is sent to indicate successful history reset.
func (b *Bot) resetResponse(u *gpt.User, evt *event.Event, msg string) error {
	u.History.ResetHistory()
	if msg != "" {
		return b.completionResponse(u, evt.RoomID, msg)
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
