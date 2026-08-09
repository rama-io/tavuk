package announce

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"tavuk/internal/commands"
	"tavuk/internal/store"
)

func (c *Command) handleSet(s *discordgo.Session, i *discordgo.InteractionCreate, sub *discordgo.ApplicationCommandInteractionDataOption) error {
	target := store.Target{GuildID: i.GuildID}
	app := ""

	for _, opt := range sub.Options {
		switch opt.Name {
		case "app":
			app = opt.StringValue()
		case "channel":
			target.ChannelID = opt.ChannelValue(s).ID
		case "role":
			target.RoleID = opt.RoleValue(s, i.GuildID).ID
		}
	}

	if app == "" || target.ChannelID == "" || target.RoleID == "" {
		return commands.Respond(s, i.Interaction, "set requires an app, a channel and a role.")
	}

	if err := c.Store.Set(app, target); err != nil {
		return err
	}

	return commands.Respond(s, i.Interaction, fmt.Sprintf(
		"Saved: **%s** will be announced to <#%s> with role <@&%s>.",
		app, target.ChannelID, target.RoleID,
	))
}
