package bot

import (
	"fmt"
	"github.com/mazzz1y/matrix-gpt/internal/gpt"
	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

// setupMessageEvent sets up the handler for incoming messages.
func (b *Bot) setupMessageEvent() {
	syncer := b.client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, b.messageHandler)
}

// setupJoinRoomEvent sets up the handler for joining rooms.
func (b *Bot) setupJoinRoomEvent() {
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

// messageHandler handles incoming messages based on their type.
func (b *Bot) messageHandler(source mautrix.EventSource, evt *event.Event) {
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

	err := b.sendResponse(user, evt)
	if err != nil {
		b.reactionResponse(evt, "❌")
		l.Err(err).Msg("response error")
	}

	user.UpdateLastMsgTime()
	l.Info().Msg("message sent")
}

// sendResponse responds to the user command.
func (b *Bot) sendResponse(user *gpt.User, evt *event.Event) (err error) {
	user.GPTMutex.Lock()
	go func() {
		b.markRead(evt)
		b.startTyping(evt.RoomID)
	}()
	defer b.stopTyping(evt.RoomID)
	defer user.GPTMutex.Unlock()

	cmd := extractCommand(evt.Content.AsMessage().Body)
	msg := trimCommand(evt.Content.AsMessage().Body)

	switch cmd {
	case HelpCommand:
		err = b.helpResponse(evt.RoomID)
	case GenerateImageCommand:
		err = b.imageResponse(evt.RoomID, msg)
	case HistoryResetCommand:
		err = b.resetResponse(user, evt, msg)
	case "":
		err = b.completionResponse(user, evt.RoomID, msg)
	default:
		err = fmt.Errorf("command \"!%s\" does not exist", cmd)
	}

	return err
}
