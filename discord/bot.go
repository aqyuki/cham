package discord

import (
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
