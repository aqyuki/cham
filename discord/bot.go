package discord

import (
	"errors"
	"fmt"

	"github.com/aqyuki/expand-bot/logging"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// Bot provides features to interact with Discord.
type Bot struct {
	token  string
	client *discordgo.Session
	logger *zap.SugaredLogger
}

// Config holds the discord configuration
type Config struct {
	Token string
}

type Option func(*Bot)

func WithLogger(logger *zap.SugaredLogger) Option {
	return func(b *Bot) {
		b.logger = logger
	}
}

// NewBot creates a new Bot instance.
func NewBot(token string, opts ...Option) Bot {
	b := Bot{
		token:  token,
		client: nil,
		logger: logging.DefaultLogger(),
	}

	for _, f := range opts {
		f(&b)
	}
	return b
}

func (b *Bot) Start() error {
	if b.client != nil {
		return errors.New("bot is already running")
	}

	session, err := discordgo.New("Bot " + b.token)
	if err != nil {
		return fmt.Errorf("failed to create session to discord because %w", err)
	}

	session.AddHandler(newMessageCreateHandler(b.logger))
	if err := session.Open(); err != nil {
		return fmt.Errorf("failed to open session to discord because %w", err)
	}

	b.client = session
	b.logger.Info("bot is running")
	return nil
}

func (b *Bot) Stop() error {
	if b.client == nil {
		return errors.New("bot is not running")
	}

	defer b.purge()
	if err := b.client.Close(); err != nil {
		return fmt.Errorf("failed to close session to discord because %w", err)
	}
	b.logger.Info("bot is stopped")
	return nil
}

// purge cleans up the bot instance.
func (b *Bot) purge() {
	b.client = nil
}
