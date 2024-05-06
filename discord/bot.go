package discord

import (
	"errors"
	"fmt"

	"github.com/aqyuki/expand-bot/logging"
	"github.com/aqyuki/expand-bot/types"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// Bot provides features to interact with Discord.
type Bot struct {
	config Config
	client *discordgo.Session
	logger *zap.SugaredLogger
}

// DiscordConfigProvider is an interface that provides the configuration for the Discord bot.
type DiscordConfigProvider interface {
	DiscordConfig() Config
}

// Config holds the discord configuration
type Config struct {
	Token types.SecretString
}

type Option func(*Bot)

func WithLogger(logger *zap.SugaredLogger) Option {
	return func(b *Bot) {
		b.logger = logger
	}
}

// NewBot creates a new Bot instance.
func NewBot(config DiscordConfigProvider, opts ...Option) Bot {
	b := Bot{
		config: config.DiscordConfig(),
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

	session, err := discordgo.New("Bot " + b.config.Token.Raw())
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
