package discord

import (
	"log"
	"os"
	"os/signal"
	"progression/league"

	"github.com/bwmarrin/discordgo"
)

type DiscordInteractionFunction func(*discordgo.Session, *discordgo.InteractionCreate)

type Bot struct {
	session         *discordgo.Session
	commands        []*discordgo.ApplicationCommand
	commandHandlers map[string]DiscordInteractionFunction
	leagueManager   *league.Manager
}

func New(token string, leagueManager *league.Manager) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		session:       session,
		leagueManager: leagueManager,
	}

	bot.commands = generateCommands()
	bot.commandHandlers = generateCommandHandlerMap(bot)

	return bot, nil
}

func generateCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Get help on how to use the bot.",
		},
		{
			Name:        "join",
			Description: "Join the upcoming league.",
		},
		{
			Name:        "start",
			Description: "Start a new league.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "set_code",
					Description: "The first set to make available to all players.",
					Required:    true,
				},
			},
		},
		{
			Name:        "redeem",
			Description: "Redeem a wild card or pack.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "card",
					Description: "Redeem a wild card to get a specific card from an already unlocked set.",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "set_code",
							Description: "The set the card belongs to.",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "collector_number",
							Description: "The collector number of the card in the given set.",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "pack",
					Description: "Redeem one or more packs to get that many packs from an already unlocked set.",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "set_code",
							Description: "The set to which the packs belong.",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "count",
							Description: "The number of packs to open.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func generateCommandHandlerMap(bot *Bot) map[string]DiscordInteractionFunction {
	commandHandlers := map[string]DiscordInteractionFunction{
		"help": bot.HelpCommand,
		"join": bot.JoinCommand,
	}
	return commandHandlers
}

func (b *Bot) Start() error {
	log.Println("Adding Ready Handler...")
	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	log.Println("Adding Interaction Handler...")
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := b.commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	err := b.session.Open()
	if err != nil {
		return err
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(b.commands))
	for i, v := range b.commands {
		cmd, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer b.session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")
	// // We need to fetch the commands, since deleting requires the command ID.
	// // We are doing this from the returned commands on line 375, because using
	// // this will delete all the commands, which might not be desirable, so we
	// // are deleting only the commands that we added.
	// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
	// if err != nil {
	// 	log.Fatalf("Could not fetch registered commands: %v", err)
	// }

	for _, v := range registeredCommands {
		err = b.session.ApplicationCommandDelete(b.session.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")

	return nil
}
