package discord

import "github.com/bwmarrin/discordgo"

func (b *Bot) HelpCommand(*discordgo.Session, *discordgo.InteractionCreate) {

}

func (b *Bot) JoinCommand(*discordgo.Session, *discordgo.InteractionCreate) {

}

func (b *Bot) ReportCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	wins := optionMap["games_won"].IntValue()
	losses := optionMap["games_lost"].IntValue()
	draws := optionMap["draws"].IntValue()
	userID := i.Member.User.ID

	err := b.leagueManager.ReportMatch(userID, int(wins), int(losses), int(draws))
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error reporting match result: " + err.Error(),
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Match result reported successfully!",
		},
	})
}
