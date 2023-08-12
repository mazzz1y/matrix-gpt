package bot

import (
	"context"
	"fmt"

	"github.com/mazzz1y/matrix-gpt/internal/gpt"
	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

// joinRoomHandler handles when the bot is invited to a room.
func (b *Bot) joinRoomHandler(source mautrix.EventSource, evt *event.Event) {
	l := log.With().
		Str("event", "join-room").
		Str("user-id", evt.Sender.String()).
		Logger()
	_, ok := b.gptClient.GetUser(evt.Sender.String())
	if ok &&
		evt.GetStateKey() == b.client.UserID.String() &&
		evt.Content.AsMember().Membership == event.MembershipInvite {
		_, err := b.client.JoinRoomByID(evt.RoomID)
		if err != nil {
			l.Err(err).Msg("join room error")
		}
	}
}

// redactionHandler handles when a previous message is redacted (deleted).
func (b *Bot) redactionHandler(source mautrix.EventSource, evt *event.Event) {
	l := log.With().
		Str("event", "redaction").
		Str("user-id", evt.Sender.String()).
		Logger()

	user, ok := b.gptClient.GetUser(evt.Sender.String())
	if !ok {
		l.Info().Msg("forbidden")
		return
	}

	reqID, ok := user.GetActiveRequestID()
	if ok && reqID == evt.Redacts.String() {
		user.CancelRequestContext(reqID)
		l.Info().Msg("message cancelled")
	}
}

// messageHandler handles incoming messages based on their type.
func (b *Bot) messageHandler(source mautrix.EventSource, evt *event.Event) {
	if b.client.UserID.String() == evt.Sender.String() {
		return
	}

	l := log.With().
		Str("event", "message").
		Str("user-id", evt.Sender.String()).
		Logger()

	user, ok := b.gptClient.GetUser(evt.Sender.String())
	if !ok {
		l.Info().Msg("forbidden")
		return
	}

	go func() {
		ctx := user.CreateRequestContext(evt.ID.String())
		defer user.CancelRequestContext(evt.ID.String())

		err := b.sendResponse(*ctx, user, evt)
		if err == context.Canceled {
			return
		}
		if err != nil {
			b.reactionResponse(evt, "❌")
			l.Err(err).Msg("response error")
		}

		user.UpdateLastMsgTime()
		l.Info().Msg("message sent")
	}()
}

// sendResponse responds to the user command.
func (b *Bot) sendResponse(ctx context.Context, user *gpt.User, evt *event.Event) (err error) {
	user.ReqMutex.Lock()
	go func() {
		b.markRead(evt)
		b.startTyping(evt.RoomID)
	}()
	defer b.stopTyping(evt.RoomID)
	defer user.ReqMutex.Unlock()

	cmd := extractCommand(evt.Content.AsMessage().Body)
	msg := trimCommand(evt.Content.AsMessage().Body)

	switch cmd {
	case HelpCommand:
		err = b.helpResponse(evt.RoomID)
	case GenerateImageCommand:
		err = b.imageResponse(ctx, evt.RoomID, msg)
	case HistoryResetCommand:
		err = b.resetResponse(ctx, user, evt, msg)
	case "":
		err = b.completionResponse(ctx, user, evt.RoomID, msg)
	default:
		err = fmt.Errorf("command \"!%s\" does not exist", cmd)
	}

	return err
}
