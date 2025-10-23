package league

import (
	"errors"
	"fmt"
	"progression/packGenerator"
	"progression/repository"
	"strconv"
)

type Manager struct {
	dataStore  repository.DataStore
	mbpgClient packGenerator.Client
}

func NewLeagueManager(dataStore repository.DataStore) *Manager {
	return &Manager{
		dataStore: dataStore,
	}
}

func (m *Manager) JoinLeague(userID string) error {
	const errMsg = "failed to join league: %w"

	player, err := m.dataStore.GetPlayer(userID)
	if err == nil {
		if player.Dropped {
			player.Dropped = false
			err = m.dataStore.UpdatePlayer(player)
			if err != nil {
				return fmt.Errorf(errMsg, err)
			}
			return nil
		}
		return fmt.Errorf(errMsg, ErrPlayerAlreadyJoined)
	}

	if !errors.Is(err, repository.ErrPlayerNotFound) {
		return fmt.Errorf(errMsg, err)
	}

	err = m.dataStore.UpdatePlayer(repository.Player{
		Id:        userID,
		WildCards: 0,
		WildPacks: 0,
	})
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}

func (m *Manager) GetBannedCards() ([]repository.Ban, error) {
	const errMsg = "failed to get banned cards: %w"

	bans, err := m.dataStore.GetBannedCards()
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return bans, nil
}

func (m *Manager) BanCard(userID, cardName string) error {
	const errMsg = "failed to ban card: %w"

	isAdmin, err := m.dataStore.IsAdmin(userID)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	if !isAdmin {
		return ErrPlayerNotAdmin
	}

	err = m.dataStore.BanCard(cardName)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}

func (m *Manager) UnbanCard(userID, cardName string) error {
	const errMsg = "failed to unban card: %w"

	isAdmin, err := m.dataStore.IsAdmin(userID)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	if !isAdmin {
		return ErrPlayerNotAdmin
	}

	err = m.dataStore.UnbanCard(cardName)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}

func (m *Manager) ReportMatch(userID string, wins, losses, draws int) error {
	const errMsg = "failed to report match: %w"

	if wins == 0 && losses == 0 && draws == 0 {
		return fmt.Errorf(errMsg, ErrInvalidMatchResult)
	}

	pairing, err := m.dataStore.GetPairing(userID)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	if isMatchReported(pairing) {
		return fmt.Errorf(errMsg, ErrMatchAlreadyReported)
	}

	if pairing.Player1 == userID {
		pairing.Wins1 = wins
		pairing.Wins2 = losses
		pairing.Draws = draws
	} else {
		pairing.Wins1 = losses
		pairing.Wins2 = wins
		pairing.Draws = draws
	}

	err = m.dataStore.UpdatePairing(pairing)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}

func isMatchReported(pairing repository.Pairing) bool {
	return pairing.Wins1 != 0 || pairing.Wins2 != 0 || pairing.Draws != 0
}

func (m *Manager) StartRound(set string) (map[string][]repository.Card, error) {
	const errMsg = "failed to start round: %w"

	players, err := m.dataStore.GetAllPlayers()
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	playerPools := make(map[string][]repository.Card)
	for _, player := range players {
		cards, err := m.mbpgClient.GetPacks(set, 10)
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}

		convertedCards := convertCardsFormat(cards)
		err = m.dataStore.StoreCards(player.Id, convertedCards)
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}

		playerPools[player.Id] = convertedCards
	}

	err = m.dataStore.StartLeague()
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return playerPools, nil
}

func convertCardsFormat(cards []packGenerator.Card) []repository.Card {
	convertedCards := make([]repository.Card, 0, len(cards))
	for _, card := range cards {
		collectorNumber, err := strconv.Atoi(card.CollectorNumber)
		if err != nil {
			// TODO there is no recovery planned from this problem
			panic(err)
		}
		convertedCards = append(convertedCards, repository.Card{
			Name:            card.Name,
			Set:             card.Set,
			CollectorNumber: collectorNumber,
			Count:           1,
		})
	}
	return convertedCards
}

func (m *Manager) GetPlayerCards(userID string) ([]repository.Card, error) {
	const errMsg = "failed to get player cards: %w"

	cards, err := m.dataStore.GetCards(userID)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return cards, nil
}

func (m *Manager) GetPlayerBalance(userID string) (repository.Player, error) {
	const errMsg = "failed to get player balance: %w"

	player, err := m.dataStore.GetPlayer(userID)
	if err != nil {
		return repository.Player{}, fmt.Errorf(errMsg, err)
	}

	return player, nil
}

func (m *Manager) DropPlayer(userID string) error {
	const errMsg = "failed to drop player: %w"

	player, err := m.dataStore.GetPlayer(userID)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	if player.Dropped {
		return ErrPlayerAlreadyDropped
	}

	pairing, err := m.dataStore.GetPairing(userID)
	if err != nil && !errors.Is(err, repository.ErrPairingNotFound) {
		return fmt.Errorf(errMsg, err)
	}

	if err == nil && !isMatchReported(pairing) {
		if pairing.Player1 == userID {
			pairing.Wins1 = 0
			pairing.Wins2 = 2
		} else {
			pairing.Wins1 = 2
			pairing.Wins2 = 0
		}
		pairing.Draws = 0

		err = m.dataStore.UpdatePairing(pairing)
		if err != nil {
			return fmt.Errorf(errMsg, err)
		}
	}

	err = m.dataStore.DropPlayer(userID)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}
