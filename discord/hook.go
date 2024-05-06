package discord

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func newMessageCreateHandler(logger *zap.SugaredLogger) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		if m.Author.Bot {
			return
		}

		var embeds []*discordgo.MessageEmbed

		// if the message contains message link, expand it
		links := extractMessageLinks(m.Content)
		for _, link := range links {
			info, err := extractMessageInfo(link)
			if err != nil {
				// TODO: add logging
				continue
			}

			// if the guild is not the same as the message, ignore it
			if info.guild != m.GuildID {
				// TODO: add logging
				continue
			}

			// if the channel is nsfw, ignore it
			citationChannel, err := s.Channel(info.channel)
			if err != nil {
				// TODO: add logging
				continue
			}
			if citationChannel.NSFW {
				// TODO: add logging
				continue
			}

			// get the message
			citationMsg, err := s.ChannelMessage(info.channel, info.message)
			if err != nil {
				// TODO: Add logging
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
				Description: citationMsg.Content,
				Timestamp:   citationMsg.Timestamp.Format(time.RFC3339),
				Footer: &discordgo.MessageEmbedFooter{
					Text: citationChannel.Name,
				},
			}

			embeds = append(embeds, embed)
		}
		// expand the message
		replyMsg := discordgo.MessageSend{
			Embeds:    embeds,
			Reference: m.Reference(),
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Users: []string{m.Author.ID},
			},
		}
		if _, err := s.ChannelMessageSendComplex(m.ChannelID, &replyMsg); err != nil {
			// TODO: add logging
			return
		}
	}
}

// func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	// Ignore all messages created by the bot itself
// 	if m.Author.Bot {
// 		return
// 	}

// 	var embeds []*discordgo.MessageEmbed

// 	// if the message contains message link, expand it
// 	links := extractMessageLinks(m.Content)
// 	for _, link := range links {
// 		info, err := extractMessageInfo(link)
// 		if err != nil {
// 			// TODO: add logging
// 			continue
// 		}

// 		// if the guild is not the same as the message, ignore it
// 		if info.guild != m.GuildID {
// 			// TODO: add logging
// 			continue
// 		}

// 		// if the channel is nsfw, ignore it
// 		citationChannel, err := s.Channel(info.channel)
// 		if err != nil {
// 			// TODO: add logging
// 			continue
// 		}
// 		if citationChannel.NSFW {
// 			// TODO: add logging
// 			continue
// 		}

// 		// get the message
// 		citationMsg, err := s.ChannelMessage(info.channel, info.message)
// 		if err != nil {
// 			// TODO: Add logging
// 			continue
// 		}

// 		// if the message has attachment(image), use it as the thumbnail of the embed
// 		var image *discordgo.MessageEmbedImage
// 		if len(citationMsg.Attachments) > 0 {
// 			image = &discordgo.MessageEmbedImage{
// 				URL: citationMsg.Attachments[0].URL,
// 			}
// 		}

// 		embed := &discordgo.MessageEmbed{
// 			Image: image,
// 			Author: &discordgo.MessageEmbedAuthor{
// 				Name:    citationMsg.Author.Username,
// 				IconURL: citationMsg.Author.AvatarURL("64"),
// 			},
// 			Description: citationMsg.Content,
// 			Timestamp:   citationMsg.Timestamp.Format(time.RFC3339),
// 			Footer: &discordgo.MessageEmbedFooter{
// 				Text: citationChannel.Name,
// 			},
// 		}

// 		embeds = append(embeds, embed)
// 	}
// 	// expand the message
// 	replyMsg := discordgo.MessageSend{
// 		Embeds:    embeds,
// 		Reference: m.Reference(),
// 		AllowedMentions: &discordgo.MessageAllowedMentions{
// 			Users: []string{m.Author.ID},
// 		},
// 	}
// 	if _, err := s.ChannelMessageSendComplex(m.ChannelID, &replyMsg); err != nil {
// 		// TODO: add logging
// 		return
// 	}
// }

var rgx = regexp.MustCompile(`https://(?:ptb\.|canary\.)?discord\.com/channels/(\d+)/(\d+)/(\d+)`)

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
	if len(segments) >= 4 {
		return message{
			guild:   segments[len(segments)-3],
			channel: segments[len(segments)-2],
			message: segments[len(segments)-1],
		}, nil
	}
	return message{}, errors.New("invalid message link")
}
