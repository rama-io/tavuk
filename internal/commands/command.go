package commands

import (
	"github.com/bwmarrin/discordgo"
)

// Definition registers it, Handle dispatches it.
type Command interface {
	Name() string
	Definition() *discordgo.ApplicationCommand
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

// Registry stores commands and dispatches interactions to them.
type Registry struct {
	byName map[string]Command
	order  []string
}

func New(cs ...Command) *Registry {
	r := &Registry{byName: map[string]Command{}}
	for _, c := range cs {
		r.byName[c.Name()] = c
		r.order = append(r.order, c.Name())
	}
	return r
}

func (r *Registry) Definitions() []*discordgo.ApplicationCommand {
	defs := make([]*discordgo.ApplicationCommand, 0, len(r.order))
	for _, name := range r.order {
		defs = append(defs, r.byName[name].Definition())
	}
	return defs
}

func (r *Registry) Names() []string {
	return append([]string(nil), r.order...)
}

// Handle routes an interaction to the matching command
func (r *Registry) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.Type {
	case discordgo.InteractionApplicationCommand, discordgo.InteractionApplicationCommandAutocomplete:
	default:
		return nil
	}
	c, ok := r.byName[i.ApplicationCommandData().Name]
	if !ok {
		return nil
	}
	return c.Handle(s, i)
}
