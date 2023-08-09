package bot

import (
	"github.com/mazzz1y/matrix-gpt/internal/gpt"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
	"time"
)

// sendAnswer responds to a user message with a GPT-based completion.
func (b *Bot) sendAnswer(u *gpt.User, evt *event.Event) error {
	if err := b.client.MarkRead(evt.RoomID, evt.ID); err != nil {
		return err
	}

	b.startTyping(evt.RoomID)
	defer b.stopTyping(evt.RoomID)

	msg := evt.Content.AsMessage().Body
	answer, err := b.gptClient.GetCompletion(u, msg)
	if err != nil {
		return err
	}

	formattedMsg := format.RenderMarkdown(answer, true, false)
	_, err = b.client.SendMessageEvent(evt.RoomID, event.EventMessage, &formattedMsg)
	return err
}

// sendReaction sends a reaction to a message.
func (b *Bot) sendReaction(evt *event.Event, emoji string) error {
	_, err := b.client.SendReaction(evt.RoomID, evt.ID, emoji)
	return err
}

// startTyping notifies the room that the bot is typing.
func (b *Bot) startTyping(roomID id.RoomID) {
	timeout := time.Duration(b.gptClient.GetTimeout()) * time.Second
	_, _ = b.client.UserTyping(roomID, true, timeout)
}

// stopTyping notifies the room that the bot has stopped typing.
func (b *Bot) stopTyping(roomID id.RoomID) {
	_, _ = b.client.UserTyping(roomID, false, 0)
}
