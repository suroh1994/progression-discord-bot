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

	_, err := m.dataStore.GetPlayer(userID)
	if err == nil {
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
