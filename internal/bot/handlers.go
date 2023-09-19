package bot

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

// joinRoomHandler handles when the bot is invited to a room.
func (b *Bot) joinRoomHandler(source mautrix.EventSource, evt *event.Event) {
	userID := evt.Sender.String()
	l := log.With().
		Str("event", "join-room").
		Str("user-id", userID).
		Logger()
	_, ok := b.users[userID]
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
	userID := evt.Sender.String()
	l := log.With().
		Str("event", "redaction").
		Str("user-id", userID).
		Logger()

	user, ok := b.users[userID]
	if !ok {
		l.Debug().Msg("forbidden")
		return
	}

	reqID, ok := user.getActiveRequestID()
	if ok && reqID == evt.Redacts.String() {
		user.cancelRequestContext(reqID)
		l.Debug().Msg("request cancelled")
	}
}

// messageHandler handles incoming messages based on their type.
func (b *Bot) messageHandler(source mautrix.EventSource, evt *event.Event) {
	userID := evt.Sender.String()
	if b.client.UserID.String() == userID {
		return
	}

	l := log.With().
		Str("event", "message").
		Str("user-id", userID).
		Logger()

	user, ok := b.users[userID]
	if !ok {
		l.Debug().Msg("forbidden")
		return
	}
	l.Debug().Msg("received request, processing")
	if user.getLastMsgTime().Add(b.historyExpire).Before(time.Now()) {
		l.Debug().Msg("history expired, resetting before processing")
		user.history.reset()
	}

	go func() {
		evtID := evt.ID.String()
		ctx := user.createRequestContext(evtID)
		defer user.cancelRequestContext(evtID)

		err := b.sendResponse(*ctx, user, evt)
		if err == context.Canceled {
			return
		}
		if err != nil {
			b.err(evt, err)
			l.Err(err).Msg("response error")
			return
		}

		user.updateLastMsgTime()
		l.Debug().Int("history-size", user.history.getSize()).Msg("response sent")
	}()
}

// sendResponse responds to the user command.
func (b *Bot) sendResponse(ctx context.Context, u *user, e *event.Event) (err error) {
	u.reqMutex.Lock()
	b.markRead(e)
	b.startTyping(e.RoomID)
	defer b.stopTyping(e.RoomID)
	defer u.reqMutex.Unlock()

	body := e.Content.AsMessage().Body
	cmd := extractCommand(body)
	msg := trimCommand(body)

	action, err := b.getAction(cmd)
	if err != nil {
		return err
	}

	return action(ctx, u, e, msg)
}

// err is a helper function to process specific error types.
func (b *Bot) err(evt *event.Event, err error) {
	switch t := err.(type) {
	case *unknownCommandError:
		b.markdownResponse(evt, true, unknownCommandMsg)
	case *openai.APIError:
		b.markdownResponse(evt, true, t.Message)
	default:
		if errors.Is(err, context.DeadlineExceeded) {
			b.markdownResponse(evt, true, timeoutMsg)
		} else {
			b.reactionResponse(evt, "❌")
		}
	}
}
