package league

import (
	"errors"
	"fmt"
	"progression/repository"
)

type Manager struct {
	dataStore repository.DataStore
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
