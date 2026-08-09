package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Respond sends an confirmation to the invoker.
func Respond(s *discordgo.Session, i *discordgo.Interaction, content string) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func LogError(i *discordgo.InteractionCreate, err error) {
	log.Printf("command %q failed: %v", i.ApplicationCommandData().Name, err)
}
