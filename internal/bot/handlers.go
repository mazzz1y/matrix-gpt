package bot

import (
	"time"

	"github.com/mazzz1y/matrix-gpt/internal/gpt"
	"github.com/mazzz1y/matrix-gpt/internal/text"
	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

// messageHandler sets up the handler for incoming messages.
func (b *Bot) messageHandler() {
	syncer := b.client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, b.msgEvtDispatcher)
}

// joinRoomHandler sets up the handler for joining rooms.
func (b *Bot) joinRoomHandler() {
	syncer := b.client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.StateMember, func(source mautrix.EventSource, evt *event.Event) {
		_, ok := b.gptClient.GetUser(evt.Sender.String())
		if ok &&
			evt.GetStateKey() == b.client.UserID.String() &&
			evt.Content.AsMember().Membership == event.MembershipInvite {
			_, err := b.client.JoinRoomByID(evt.RoomID)
			if err != nil {
				return
			}
		}
	})
}

// historyResetHandler checks for the reset command and resets history if found.
func (b *Bot) historyResetHandler(user *gpt.User, evt *event.Event) (ok bool) {
	if text.HasPrefixIgnoreCase(evt.Content.AsMessage().Body, "!reset") {
		user.History.ResetHistory()
		_ = b.sendReaction(evt, "✅")
		return true
	}
	return false
}

// historyExpireHandler checks if the history for a user has expired and resets if necessary.
func (b *Bot) historyExpireHandler(user *gpt.User) {
	if user.LastMsg.Add(time.Duration(b.historyExpire) * time.Hour).Before(time.Now()) {
		user.History.ResetHistory()
	}
}

// msgEvtDispatcher dispatches incoming messages to their appropriate handlers.
func (b *Bot) msgEvtDispatcher(source mautrix.EventSource, evt *event.Event) {
	// Ignore messages sent by the bot itself
	if b.client.UserID.String() == evt.Sender.String() {
		return
	}

	l := log.With().
		Str("component", "handler").
		Str("user_id", evt.Sender.String()).
		Logger()

	user, ok := b.gptClient.GetUser(evt.Sender.String())
	if !ok {
		l.Info().Msg("forbidden")
		return
	}

	if b.historyResetHandler(user, evt) {
		l.Info().Msg("reset history")
		return
	}

	err := b.sendAnswer(user, evt)
	if err != nil {
		l.Err(err).Msg("failed to send message")
	} else {
		l.Info().Msg("sending answer")
	}
}
