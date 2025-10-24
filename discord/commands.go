package discord

import (
	"errors"
	"fmt"
	"progression/league"
	"progression/repository"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) HelpCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	message := "Please check the command section in the bot README here: https://github.com/suroh1994/progression-discord-bot?tab=readme-ov-file#commands"

	return b.SendMessage(s, i, message)
}

func (b *Bot) SetsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var message string
	sets, err := b.leagueManager.GetSets()
	if err != nil {
		message = "Error getting sets: " + err.Error()
	} else {
		if len(sets) == 0 {
			message = "There are no unlocked sets."
		} else {
			var builder strings.Builder
			builder.WriteString("```\n")
			for _, set := range sets {
				builder.WriteString(fmt.Sprintf("%s\n", set.SetCode))
			}
			builder.WriteString("```")
			message = builder.String()
		}
	}

	return b.SendMessage(s, i, message)
}

func (b *Bot) BansCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var message string
	bans, err := b.leagueManager.GetBannedCards()
	if err != nil {
		message = "Error getting banned cards: " + err.Error()
	} else {
		if len(bans) == 0 {
			message = "There are no banned cards."
		} else {
			var builder strings.Builder
			builder.WriteString("```\n")
			for _, ban := range bans {
				builder.WriteString(fmt.Sprintf("%s\n", ban.CardName))
			}
			builder.WriteString("```")
			message = builder.String()
		}
	}

	return b.SendMessage(s, i, message)
}

func (b *Bot) BanCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	userID := i.Member.User.ID
	commandData := i.ApplicationCommandData()
	cardName := commandData.Options[0].StringValue()

	var message string
	err := b.leagueManager.BanCard(userID, cardName)
	if err != nil {
		switch {
		case errors.Is(err, league.ErrPlayerNotAdmin):
			message = "You are not an admin."
		default:
			message = "Error banning card: " + err.Error()
		}
	} else {
		message = fmt.Sprintf("Banned %s.", cardName)
	}

	return b.SendMessage(s, i, message)
}

func (b *Bot) UnbanCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	userID := i.Member.User.ID
	commandData := i.ApplicationCommandData()
	cardName := commandData.Options[0].StringValue()

	var message string
	err := b.leagueManager.UnbanCard(userID, cardName)
	if err != nil {
		switch {
		case errors.Is(err, league.ErrPlayerNotAdmin):
			message = "You are not an admin."
		default:
			message = "Error unbanning card: " + err.Error()
		}
	} else {
		message = fmt.Sprintf("Unbanned %s.", cardName)
	}

	return b.SendMessage(s, i, message)
}

func (b *Bot) DropCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	userID := i.Member.User.ID

	var message string
	err := b.leagueManager.DropPlayer(userID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrPlayerNotFound):
			message = "You are not part of the current league."
		case errors.Is(err, league.ErrPlayerAlreadyDropped):
			message = "You have already dropped from the league."
		default:
			message = "Error dropping from the league: " + err.Error()
		}
	} else {
		message = "You have been successfully removed from the league."
	}

	return b.SendMessage(s, i, message)
}

func (b *Bot) BalanceCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	userID := i.Member.User.ID
	var message string
	player, err := b.leagueManager.GetPlayerBalance(userID)
	if err != nil {
		message = "Error getting your balance: " + err.Error()
	} else {
		message = fmt.Sprintf("Wild cards: %d\nWild packs: %d", player.WildCards, player.WildPacks)
	}

	return b.SendMessage(s, i, message)
}

func (b *Bot) PoolCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	userID := i.Member.User.ID
	var message string
	cards, err := b.leagueManager.GetPlayerCards(userID)
	if err != nil {
		message = "Error getting your card pool: " + err.Error()
	} else {
		message = formatCardList(cards)
	}

	return b.SendMessage(s, i, message)
}

func formatCardList(cards []repository.Card) string {
	if len(cards) == 0 {
		return "You currently have no cards in your pool."
	}
	var builder strings.Builder
	builder.WriteString("```\n")
	for _, card := range cards {
		builder.WriteString(fmt.Sprintf("%d %s\n", card.Count, card.Name))
	}
	builder.WriteString("```")
	return builder.String()
}

func (b *Bot) JoinCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	userID := i.Member.User.ID

	var message string
	err := b.leagueManager.JoinLeague(userID)
	if err != nil {
		switch {
		case errors.Is(err, league.ErrPlayerAlreadyJoined):
			message = "You've already joined the league."
		default:
			message = "Error joining the league: " + err.Error()
		}
	} else {
		message = "You've joined the league."
	}

	return b.SendMessage(s, i, message)
}

func (b *Bot) ReportCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	userID := i.Member.User.ID
	commandData := i.ApplicationCommandData()
	wins := commandData.GetOption("games_won").IntValue()
	losses := commandData.GetOption("games_lost").IntValue()
	draws := commandData.GetOption("draws").IntValue()

	var message string
	err := b.leagueManager.ReportMatch(userID, int(wins), int(losses), int(draws))
	if err != nil {
		switch {
		case errors.Is(err, league.ErrInvalidMatchResult):
			message = fmt.Sprintf("The given match result is invalid. Given: %d wins, %d losses and %d draws.", wins, losses, int(draws))
		case errors.Is(err, league.ErrMatchAlreadyReported):
			message = "Your match has already been reported."
		default:
			message = "Error reporting match result: " + err.Error()
		}
	} else {
		message = "Match result reported successfully!"
	}

	return b.SendMessage(s, i, message)
}
