package commands

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

type stubCommand struct {
	name    string
	handled bool
}

func (c *stubCommand) Name() string { return c.name }

func (c *stubCommand) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{Name: c.name}
}

func (c *stubCommand) Handle(*discordgo.Session, *discordgo.InteractionCreate) error {
	c.handled = true
	return nil
}

func TestRegistryDispatchesCommandAndAutocomplete(t *testing.T) {
	for _, typ := range []discordgo.InteractionType{
		discordgo.InteractionApplicationCommand,
		discordgo.InteractionApplicationCommandAutocomplete,
	} {
		cmd := &stubCommand{name: "announce"}
		r := New(cmd)
		if err := r.Handle(nil, &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Type: typ,
				Data: discordgo.ApplicationCommandInteractionData{Name: "announce"},
			},
		}); err != nil {
			t.Fatalf("Handle(%v): %v", typ, err)
		}
		if !cmd.handled {
			t.Fatalf("interaction of type %v was not dispatched", typ)
		}
	}
}

func TestRegistryIgnoresUnknownAndNonCommandInteractions(t *testing.T) {
	cmd := &stubCommand{name: "announce"}
	r := New(cmd)
	if err := r.Handle(nil, &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionMessageComponent,
			Data: discordgo.MessageComponentInteractionData{},
		},
	}); err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if cmd.handled {
		t.Fatal("unknown interaction must not be dispatched")
	}
}
