package main

import (
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/iopred/bruxism"
	"github.com/iopred/bruxism/carbonitexplugin"
	"github.com/iopred/bruxism/inviteplugin"
	"github.com/iopred/bruxism/playingplugin"
	"github.com/iopred/bruxism/statsplugin"
	"github.com/iopred/keystone/keystoneplugin"
)

var discordToken string
var discordEmail string
var discordPassword string
var discordApplicationClientID string
var discordOwnerUserID string
var carbonitexKey string

func init() {
	flag.StringVar(&discordToken, "discordtoken", "", "Discord token.")
	flag.StringVar(&discordEmail, "discordemail", "", "Discord account email.")
	flag.StringVar(&discordPassword, "discordpassword", "", "Discord account password.")
	flag.StringVar(&discordOwnerUserID, "discordowneruserid", "", "Discord owner user id.")
	flag.StringVar(&discordApplicationClientID, "discordapplicationclientid", "", "Discord application client id.")
	flag.StringVar(&carbonitexKey, "carbonitexkey", "", "Carbonitex key for discord server count tracking.")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
}

func main() {
	q := make(chan bool)

	// Set our variables.
	bot := bruxism.NewBot()

	// Generally CommandPlugins don't hold state, so we share one instance of the command plugin for all services.
	cp := bruxism.NewCommandPlugin()
	cp.AddCommand("invite", inviteplugin.InviteCommand, inviteplugin.InviteHelp)
	cp.AddCommand("join", inviteplugin.InviteCommand, nil)
	cp.AddCommand("stats", statsplugin.StatsCommand, statsplugin.StatsHelp)
	cp.AddCommand("info", statsplugin.StatsCommand, nil)
	cp.AddCommand("stat", statsplugin.StatsCommand, nil)
	cp.AddCommand("quit", func(bot *bruxism.Bot, service bruxism.Service, message bruxism.Message, args string, parts []string) {
		if service.IsBotOwner(message) {
			q <- true
		}
	}, nil)

	// Register the Discord service if we have an email or token.
	if (discordEmail != "" && discordPassword != "") || discordToken != "" {
		var discord *bruxism.Discord
		if discordToken != "" {
			discord = bruxism.NewDiscord(discordToken)
		} else {
			discord = bruxism.NewDiscord(discordEmail, discordPassword)
		}
		discord.ApplicationClientID = discordApplicationClientID
		discord.OwnerUserID = discordOwnerUserID
		bot.RegisterService(discord)

		bot.RegisterPlugin(discord, cp)
		bot.RegisterPlugin(discord, playingplugin.New())
		bot.RegisterPlugin(discord, keystoneplugin.New())
		if carbonitexKey != "" {
			bot.RegisterPlugin(discord, carbonitexplugin.New(carbonitexKey))
		}
	}

	// Start all our services.
	bot.Open()

	// Wait for a termination signal, while saving the bot state every minute. Save on close.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	t := time.Tick(1 * time.Minute)

out:
	for {
		select {
		case <-q:
			break out
		case <-c:
			break out
		case <-t:
			bot.Save()
		}
	}

	bot.Save()
}
