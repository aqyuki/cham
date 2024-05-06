package config

import (
	"os"

	"github.com/aqyuki/expand-bot/discord"
	"github.com/aqyuki/expand-bot/types"
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
			Token: types.SecretString(os.Getenv("DISCORD_TOKEN")),
		},
	}
}
