package announce

import (
	"github.com/bwmarrin/discordgo"

	"tavuk/internal/commands"
	"tavuk/internal/store"
)

type Command struct {
	Store *store.Store
}

func (c *Command) Name() string { return "announce" }

func (c *Command) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:                     c.Name(),
		Description:              "Announce an app release",
		DefaultMemberPermissions: memberPermissions(),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "set",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Set the target channel and role for an app",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "app",
						Description: "App name, e.g. \"Tui\"",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "channel",
						Description: "Channel to post announcements to",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "Role to mention in announcements",
						Required:    true,
					},
				},
			},
			{
				Name:        "publish",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Publish a release announcement for an app",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:         "app",
						Type:         discordgo.ApplicationCommandOptionString,
						Description:  "App name configured with /announce set",
						Required:     true,
						Autocomplete: true,
					},
					{
						Name:        "version",
						Type:        discordgo.ApplicationCommandOptionString,
						Description: "Release version, e.g. \"10\"",
						Required:    true,
					},
					{
						Name:        "changes",
						Type:        discordgo.ApplicationCommandOptionString,
						Description: "Feature list, e.g. \"- feat 1 - feat 2\"",
						Required:    true,
					},
				},
			},
		},
	}
}

// Only admin.
func memberPermissions() *int64 {
	perm := int64(discordgo.PermissionAdministrator)
	return &perm
}

// Handle dispatches /announce subcommands and autocomplete requests.
func (c *Command) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
		return c.handleAutocomplete(s, i)
	}
	if i.Type != discordgo.InteractionApplicationCommand {
		return nil
	}

	data := i.ApplicationCommandData()
	if len(data.Options) == 0 {
		return nil
	}

	sub := data.Options[0]
	var err error
	switch sub.Name {
	case "set":
		err = c.handleSet(s, i, sub)
	case "publish":
		err = c.handlePublish(s, i, sub)
	}

	if err != nil {
		commands.LogError(i, err)
	}
	return nil
}
