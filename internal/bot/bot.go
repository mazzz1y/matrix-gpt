package bot

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mazzz1y/matrix-gpt/internal/gpt"
	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/cryptohelper"
	"maunium.net/go/mautrix/event"
)

type Bot struct {
	client        *mautrix.Client
	gptClient     *gpt.Gpt
	selfProfile   mautrix.RespUserProfile
	historyExpire time.Duration
	users         map[string]*user
	actions       map[string]action
}

// NewBot initializes a new Matrix bot instance.
func NewBot(serverUrl, userID, password, sqlitePath string, historyExpire, historyLimit int, userIDs []string, gpt *gpt.Gpt) (*Bot, error) {
	client, err := mautrix.NewClient(serverUrl, "", "")
	if err != nil {
		return nil, err
	}

	crypto, err := cryptohelper.NewCryptoHelper(client, []byte("1337"), sqlitePath)
	if err != nil {
		return nil, err
	}

	crypto.LoginAs = &mautrix.ReqLogin{
		Type:       mautrix.AuthTypePassword,
		Identifier: mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: userID},
		Password:   password,
	}

	if err := crypto.Init(); err != nil {
		return nil, err
	}

	client.Crypto = crypto
	profile, err := client.GetProfile(client.UserID)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("matrix-username", profile.DisplayName).
		Str("gpt-model", gpt.GetModel()).
		Float64("gpt-timeout", gpt.GetTimeout().Seconds()).
		Int("history-limit", historyLimit).
		Int("history-expire", historyExpire).
		Msg("connected to matrix")

	users := make(map[string]*user)
	for _, id := range userIDs {
		users[id] = newGptUser(historyLimit)
	}

	return &Bot{
		client:        client,
		gptClient:     gpt,
		selfProfile:   *profile,
		users:         users,
		historyExpire: time.Duration(historyExpire) * time.Hour,
	}, nil
}

// StartHandler initializes bot event handlers and starts the matrix client sync.
func (b *Bot) StartHandler() error {
	b.initBotActions()

	syncer := b.client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, b.messageHandler)
	syncer.OnEventType(event.EventRedaction, b.redactionHandler)
	syncer.OnEventType(event.StateMember, b.joinRoomHandler)
	return b.client.Sync()
}
