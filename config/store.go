package config

import (
	"os"

	"github.com/aqyuki/expand-bot/discord"
)

// Store manages the configuration of the application.
type Store struct {
	discord discord.Config
}

func (s Store) DiscordConfig() discord.Config {
	return s.discord
}

func NewStore() Store {
	return Store{
		discord: discord.Config{
			Token: os.Getenv("DISCORD_TOKEN"),
		},
	}
}
