package bot

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/mazzz1y/matrix-gpt/internal/gpt"
	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/cryptohelper"
	"time"
)

type Bot struct {
	client        *mautrix.Client
	gptClient     *gpt.Gpt
	selfProfile   mautrix.RespUserProfile
	historyExpire time.Duration
}

// NewBot initializes a new Matrix bot instance.
func NewBot(serverUrl, userID, password, sqlitePath string, historyExpire int, gpt *gpt.Gpt) (*Bot, error) {
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
		Int("gpt-timeout", gpt.GetTimeout()).
		Int("history-limit", gpt.GetHistoryLimit()).
		Int("history-expire", historyExpire).
		Msg("connected to matrix")

	return &Bot{
		client:        client,
		gptClient:     gpt,
		selfProfile:   *profile,
		historyExpire: time.Duration(historyExpire) * time.Hour,
	}, nil
}

// StartHandler initializes bot event handlers and starts the matrix client sync.
func (b *Bot) StartHandler() error {
	logger := log.With().Str("component", "handler").Logger()

	b.setupJoinRoomEvent()
	b.setupMessageEvent()

	logger.Info().Msg("started handler")
	return b.client.Sync()
}
