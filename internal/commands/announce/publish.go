package announce

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// this acknowledges the interaction first so slow calls don't hit Discord's response timeout, then posts the announcement.
func (c *Command) handlePublish(s *discordgo.Session, i *discordgo.InteractionCreate, sub *discordgo.ApplicationCommandInteractionDataOption) error {
	var app, version, changes string
	for _, opt := range sub.Options {
		switch opt.Name {
		case "app":
			app = opt.StringValue()
		case "version":
			version = opt.StringValue()
		case "changes":
			changes = opt.StringValue()
		}
	}

	if err := deferInteraction(s, i.Interaction); err != nil {
		return err
	}

	var status string
	switch {
	case app == "" || version == "" || changes == "":
		status = "publish requires an app, a version and a list of changes."
	default:
		target, ok := c.Store.Get(app)
		switch {
		case !ok:
			status = fmt.Sprintf("No target configured for **%s**. Set one first with `/announce set`.", app)
		default:
			msg, err := s.ChannelMessageSend(target.ChannelID, buildAnnouncement(app, version, changes, target.RoleID))
			if err != nil {
				_ = reply(s, i.Interaction, "The announcement could not be sent. Check the bot's permissions on the target channel.")
				return err
			}
			status = fmt.Sprintf("Announced **%s** to <#%s>. (id: %s)",
				releaseTitle(app, version), msg.ChannelID, msg.ID)
		}
	}

	return reply(s, i.Interaction, status)
}

// this acknowledges the interaction; the reply comes later.
func deferInteraction(s *discordgo.Session, i *discordgo.Interaction) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Flags: discordgo.MessageFlagsEphemeral},
	})
}

// this sets the interaction's final status.
func reply(s *discordgo.Session, i *discordgo.Interaction, status string) error {
	_, err := s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Content: &status})
	return err
}
