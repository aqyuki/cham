package discord

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aqyuki/cham/logging"
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

	session.AddHandler(b.expandMessageLink)
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
	if err := b.client.Close(); err != nil {
		return fmt.Errorf("failed to close session to discord because %w", err)
	}
	b.logger.Info("bot is stopped")
	return nil
}

func (b *Bot) expandMessageLink(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.Bot {
		b.logger.Info("skip message because it was created by the bot itself")
		return
	}

	// if the message contains message link, expand it
	links := extractMessageLinks(m.Content)
	if len(links) == 0 {
		b.logger.Info("skip message because it does not contain message links")
		return
	}

	var embeds []*discordgo.MessageEmbed

	for _, link := range links {
		info, err := extractMessageInfo(link)
		if err != nil {
			b.logger.Error("failed to extract message info", zap.Error(err))
			continue
		}

		// if the guild is not the same as the message, ignore it
		if info.guild != m.GuildID {
			b.logger.Info("skip message because the guild is not the same as the message")
			continue
		}

		// if the channel is nsfw, ignore it
		citationChannel, err := s.Channel(info.channel)
		if err != nil {
			b.logger.Error("failed to get channel", zap.Error(err))
			continue
		}
		if citationChannel.NSFW {
			b.logger.Info("skip message because the channel is nsfw")
			continue
		}

		// get the message
		citationMsg, err := s.ChannelMessage(info.channel, info.message)
		if err != nil {
			b.logger.Error("failed to get message", zap.Error(err))
			continue
		}

		// if the message has attachment(image), use it as the thumbnail of the embed
		var image *discordgo.MessageEmbedImage
		if len(citationMsg.Attachments) > 0 {
			image = &discordgo.MessageEmbedImage{
				URL: citationMsg.Attachments[0].URL,
			}
		}

		embed := &discordgo.MessageEmbed{
			Image: image,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    citationMsg.Author.Username,
				IconURL: citationMsg.Author.AvatarURL("64"),
			},
			Color:       0x7fffff,
			Description: citationMsg.Content,
			Timestamp:   citationMsg.Timestamp.Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: citationChannel.Name,
			},
		}

		embeds = append(embeds, embed)
		b.logger.Debug("expanded message", zap.Any("embed", embed))
	}
	// expand the message
	replyMsg := discordgo.MessageSend{
		Embeds:    embeds,
		Reference: m.Reference(),
		AllowedMentions: &discordgo.MessageAllowedMentions{
			RepliedUser: true,
		},
	}
	if _, err := s.ChannelMessageSendComplex(m.ChannelID, &replyMsg); err != nil {
		b.logger.Error("failed to send message", zap.Error(err))
		return
	}
	b.logger.Info("expanded message", zap.String("channel", m.ChannelID), zap.String("message", m.ID))
}

var rgx = regexp.MustCompile(`https://(?:ptb\.|canary\.)?discord(app)?\.com/channels/(\d+)/(\d+)/(\d+)`)

func extractMessageLinks(s string) []string {
	return rgx.FindAllString(s, -1)
}

type message struct {
	guild   string
	channel string
	message string
}

// extractMessageInfo extracts the channel ID and message ID from the message link.
func extractMessageInfo(link string) (info message, err error) {
	segments := strings.Split(link, "/")
	if len(segments) < 4 {
		return message{}, errors.New("invalid message link")
	}
	return message{
		guild:   segments[len(segments)-3],
		channel: segments[len(segments)-2],
		message: segments[len(segments)-1],
	}, nil
}
