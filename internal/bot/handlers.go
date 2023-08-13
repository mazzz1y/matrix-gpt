package bot

import (
	"context"

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
		if _, ok := err.(*unknownCommandError); ok {
			b.reactionResponse(evt, "❌")
			return
		}
		if err != nil {
			b.reactionResponse(evt, "❌")
			l.Err(err).Msg("response error")
			return
		}

		user.UpdateLastMsgTime()
		l.Info().Msg("message sent")
	}()
}

// sendResponse responds to the user command.
func (b *Bot) sendResponse(ctx context.Context, user *gpt.User, evt *event.Event) (err error) {
	user.ReqMutex.Lock()
	b.markRead(evt)
	b.startTyping(evt.RoomID)
	defer b.stopTyping(evt.RoomID)
	defer user.ReqMutex.Unlock()

	inputCmd := extractCommand(evt.Content.AsMessage().Body)
	msg := trimCommand(evt.Content.AsMessage().Body)

	switch {
	case commandIs(inputCmd, helpCommand):
		err = b.helpResponse(evt.RoomID)
	case commandIs(inputCmd, generateImageCommand):
		err = b.imageResponse(ctx, evt.RoomID, msg)
	case commandIs(inputCmd, historyResetCommand):
		err = b.resetResponse(ctx, user, evt, msg)
	case commandIs(inputCmd, emptyCommand):
		err = b.completionResponse(ctx, user, evt.RoomID, msg)
	default:
		err = &unknownCommandError{cmd: inputCmd}
	}

	return err
}
