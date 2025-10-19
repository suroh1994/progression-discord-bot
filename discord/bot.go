package discord

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"progression/league"

	"github.com/bwmarrin/discordgo"
)

type InteractionFunction func(*discordgo.Session, *discordgo.InteractionCreate)

type Bot struct {
	session         *discordgo.Session
	commands        []*discordgo.ApplicationCommand
	commandHandlers map[string]InteractionFunction
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
			Name:        "pool",
			Description: "Get a list of all cards in your card pool.",
		},
		{
			Name:        "report",
			Description: "Report match results.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "games_won",
					Description: "The number of games in the match won by the reporting player.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "games_lost",
					Description: "The number of games in the match won by the opponent of the reporting player.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "draws",
					Description: "The number of games in the match ending in a draw.",
					Required:    true,
				},
			},
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

func generateCommandHandlerMap(bot *Bot) map[string]InteractionFunction {
	commandHandlers := map[string]InteractionFunction{
		"help":   WithErrorLogging(bot.HelpCommand),
		"join":   WithErrorLogging(bot.JoinCommand),
		"pool":   WithErrorLogging(bot.PoolCommand),
		"report": WithErrorLogging(bot.ReportCommand),
	}
	return commandHandlers
}

func (b *Bot) Start() error {
	slog.Info("Adding Ready Handler...")
	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		slog.Info("Successfully logged in.",
			"botname", fmt.Sprintf("%s#%s", s.State.User.Username, s.State.User.Discriminator),
		)
	})

	slog.Info("Adding Interaction Handler...")
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := b.commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	err := b.session.Open()
	if err != nil {
		return err
	}

	slog.Info("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(b.commands))
	for i, v := range b.commands {
		cmd, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", v)
		if err != nil {
			slog.Error("Cannot create command", "command", v.Name, "error", err)
			return err
		}
		registeredCommands[i] = cmd
	}

	defer b.session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	slog.Info("Press Ctrl+C to exit")
	<-stop

	slog.Info("Removing commands...")
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
			slog.Error("Cannot delete command", "command", v.Name, "error", err)
			return err
		}
	}

	slog.Info("Gracefully shutting down.")

	return nil
}

func (b *Bot) SendMessage(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

func WithErrorLogging(f func(*discordgo.Session, *discordgo.InteractionCreate) error) InteractionFunction {
	wrapFunc := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		userID := i.Member.User.ID
		err := f(s, i)
		if err != nil {
			slog.Error("failed to report error to user", "error", err, "user", userID)
		}
	}
	return wrapFunc
}
