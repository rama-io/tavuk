package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"

	"tavuk/internal/commands"
	"tavuk/internal/commands/announce"
	"tavuk/internal/config"
	"tavuk/internal/store"
)

type Bot struct {
	Session  *discordgo.Session
	Commands *commands.Registry
	Cfg      *config.Config
}

func New(s *discordgo.Session, st *store.Store, cfg *config.Config) *Bot {
	registry := commands.New(&announce.Command{Store: st})
	return &Bot{Session: s, Commands: registry, Cfg: cfg}
}

// RegisterCommands overwrites the bot's commands. If a guild ID is configured, the commands are registered only in that guild; otherwise they are registered globally.
func (b *Bot) RegisterCommands() ([]*discordgo.ApplicationCommand, error) {
	appID := b.Session.State.User.ID
	guildID := ""
	if b.Cfg != nil {
		guildID = b.Cfg.GuildID
	}

	commands, err := b.Session.ApplicationCommandBulkOverwrite(appID, guildID, b.Commands.Definitions())
	if err != nil {
		return nil, err
	}

	if guildID != "" {
		if err := b.clearGlobalCommands(appID); err != nil {
			log.Printf("warning: could not clear global copy of commands: %v", err)
		}
	}
	return commands, nil
}

func (b *Bot) clearGlobalCommands(appID string) error {
	registered, err := b.Session.ApplicationCommands(appID, "")
	if err != nil {
		return err
	}
	for _, name := range b.Commands.Names() {
		for _, cmd := range registered {
			if cmd.Name == name {
				if err := b.Session.ApplicationCommandDelete(appID, "", cmd.ID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (b *Bot) OnReady(*discordgo.Session, *discordgo.Ready) {
	_ = b.Session.UpdateGameStatus(0, "Tavuk is watching you")
}

func (b *Bot) OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = b.Commands.Handle(s, i)
}
