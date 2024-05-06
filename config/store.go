package config

import "github.com/aqyuki/expand-bot/discord"

// Store manages the configuration of the application.
type Store struct {
	discord discord.Config
}

func (s Store) DiscordConfig() discord.Config {
	return s.discord
}
