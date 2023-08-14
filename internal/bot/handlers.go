package bot

import (
	"context"
	"time"

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
	_, ok := b.getUser(evt.Sender.String())
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

	user, ok := b.getUser(evt.Sender.String())
	if !ok {
		l.Info().Msg("forbidden")
		return
	}

	reqID, ok := user.getActiveRequestID()
	if ok && reqID == evt.Redacts.String() {
		user.cancelRequestContext(reqID)
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

	user, ok := b.getUser(evt.Sender.String())
	if !ok {
		l.Info().Msg("forbidden")
		return
	}

	if user.getLastMsgTime().Add(b.historyExpire).Before(time.Now()) {
		user.history.reset()
		l.Info().Msg("history expired, resetting")
	}

	go func() {
		ctx := user.createRequestContext(evt.ID.String())
		defer user.cancelRequestContext(evt.ID.String())

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

		user.updateLastMsgTime()
		l.Info().Msg("message sent")
	}()
}

// sendResponse responds to the user command.
func (b *Bot) sendResponse(ctx context.Context, u *user, e *event.Event) (err error) {
	u.reqMutex.Lock()
	b.markRead(e)
	b.startTyping(e.RoomID)
	defer b.stopTyping(e.RoomID)
	defer u.reqMutex.Unlock()

	cmd := extractCommand(e.Content.AsMessage().Body)
	msg := trimCommand(e.Content.AsMessage().Body)

	action, err := b.getAction(cmd)
	if err != nil {
		return err
	}

	return action(ctx, u, e, msg)
}

// getUser retrieves the User instance associated with the given ID.
func (b *Bot) getUser(id string) (u *user, ok bool) {
	u, ok = b.users[id]
	return
}
