package discord

import (
	"errors"
	"fmt"

	"github.com/aqyuki/expand-bot/types"
	"github.com/bwmarrin/discordgo"
)

// Bot provides features to interact with Discord.
type Bot struct {
	config Config
	client *discordgo.Session
}

// DiscordConfigProvider is an interface that provides the configuration for the Discord bot.
type DiscordConfigProvider interface {
	DiscordConfig() Config
}

// Config holds the discord configuration
type Config struct {
	Token types.SecretString
}

// NewBot creates a new Bot instance.
func NewBot(config DiscordConfigProvider) Bot {
	b := Bot{
		config: config.DiscordConfig(),
		client: nil,
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

	if err := b.client.Open(); err != nil {
		return fmt.Errorf("failed to open session to discord because %w", err)
	}
	b.client = session
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
	return nil
}

// purge cleans up the bot instance.
func (b *Bot) purge() {
	b.client = nil
}
