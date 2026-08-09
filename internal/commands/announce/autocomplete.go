package announce

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (c *Command) handleAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	if len(data.Options) == 0 {
		return nil
	}

	sub := data.Options[0]
	for _, opt := range sub.Options {
		if !opt.Focused {
			continue
		}

		query := strings.ToLower(strings.TrimSpace(opt.StringValue()))
		choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, 25)
		for _, app := range c.Store.Apps() {
			if len(choices) == 25 {
				break
			}
			if query != "" && !strings.Contains(strings.ToLower(app), query) {
				continue
			}
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{Name: app, Value: app})
		}

		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{Choices: choices},
		})
	}
	return nil
}
