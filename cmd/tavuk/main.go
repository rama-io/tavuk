package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"tavuk/internal/bot"
	"tavuk/internal/config"
	"tavuk/internal/store"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file found (%v); using environment variables", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	st, err := store.Open(filepath.Join(cfg.DataDir, "apps.json"))
	if err != nil {
		log.Fatalf("store: %v", err)
	}

	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatalf("discord session: %v", err)
	}

	app := bot.New(session, st, cfg)
	session.AddHandler(app.OnReady)
	session.AddHandler(app.OnInteractionCreate)
	session.AddHandler(app.OnMessageCreate)

	if err := session.Open(); err != nil {
		log.Fatalf("open gateway: %v", err)
	}
	defer session.Close()

	if session.State.User != nil {
		log.Printf("Logged in as %s (id %s)", session.State.User.Username, session.State.User.ID)
	}

	inConfiguredGuild := false
	for _, g := range session.State.Guilds {
		log.Printf("In guild: %s (id %s)", g.Name, g.ID)
		if cfg.GuildID != "" && g.ID == cfg.GuildID {
			inConfiguredGuild = true
		}
	}
	if cfg.GuildID != "" && !inConfiguredGuild {
		log.Printf("WARNING: bot is not in the guild configured in GUILD_ID (%s) — commands will not appear there", cfg.GuildID)
	}

	registered, err := app.RegisterCommands()
	if err != nil {
		log.Fatalf("register commands: %v", err)
	}
	scope := "globally"
	if cfg.GuildID != "" {
		scope = "to guild " + cfg.GuildID
	}
	for _, c := range registered {
		log.Printf("Registered command %q (id %s) %s", c.Name, c.ID, scope)
	}

	log.Println("Tavuk is up. If the command is missing in Discord, re-invite the bot with the applications.commands scope.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutting down")
}
