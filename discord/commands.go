package discord

import (
	"errors"
	"fmt"
	"log/slog"
	"progression/league"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) HelpCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Please check the command section in the bot README here: https://github.com/suroh1994/progression-discord-bot?tab=readme-ov-file#commands",
		},
	})
	if err != nil {
		slog.Warn("failed to print help command", "error", err)
	}
}

func (b *Bot) JoinCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	const responseMsg = "You've successfully joined the league, %s!"

	displayName := i.Member.DisplayName()
	userID := i.Member.User.ID

	response := fmt.Sprintf(responseMsg, displayName)
	err := b.leagueManager.JoinLeague(userID)
	if err != nil {
		slog.Warn("failed to join league", "error", err)

		if errors.Is(err, league.ErrPlayerAlreadyJoined) {
			response = fmt.Sprintf("You're already in this league, %s.", displayName)
		} else {
			response = fmt.Sprintf("There was an error joining the league, please try again later.")
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		slog.Warn("failed to print join response", "error", err)
	}
}

func (b *Bot) StartCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Please check the command section in the bot README here: https://github.com/suroh1994/progression-discord-bot?tab=readme-ov-file#commands",
		},
	})
	if err != nil {
		slog.Warn("failed to print start response", "error", err)
	}

}
