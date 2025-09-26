package repository

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// IMPORTANT: (re-)start the database with `make run-pgdb` before you run these tests

func TestInsertCardPool(t *testing.T) {
	dataStore := New("localhost", 5432, "postgres", "postgres", "progression")
	err := dataStore.Connect()
	assert.NoError(t, err)

	cards := []Card{
		{
			Set:             "IKO",
			CollectorNumber: 1,
		},
	}
	playerID := "test_player" + strconv.FormatInt(time.Now().UnixMilli(), 10)

	err = dataStore.StoreCards(playerID, cards)
	assert.NoError(t, err, "failed to store cards")
}

func TestGetCardPool(t *testing.T) {
	dataStore := New("localhost", 5432, "postgres", "postgres", "progression")
	err := dataStore.Connect()
	assert.NoError(t, err)

	cards := []Card{
		{
			Set:             "IKO",
			CollectorNumber: 1,
		},
		{
			Set:             "IKO",
			CollectorNumber: 2,
		},
		{
			Set:             "IKO",
			CollectorNumber: 1,
		},
	}
	playerID := "test_player" + strconv.FormatInt(time.Now().UnixMilli(), 10)

	err = dataStore.StoreCards(playerID, cards)
	assert.NoError(t, err, "failed to store cards")

	storedCards, err := dataStore.GetCards(playerID)
	assert.NoError(t, err, "failed to get cards")

	assert.Len(t, storedCards, 2, "expected 2 different cards")
}

func TestCardPoolDeduplicate(t *testing.T) {
	dataStore := New("localhost", 5432, "postgres", "postgres", "progression")
	err := dataStore.Connect()
	assert.NoError(t, err)

	cards := []Card{
		{
			Set:             "IKO",
			CollectorNumber: 1,
		},
		{
			Set:             "IKO",
			CollectorNumber: 1,
		},
	}
	playerID := "test_player" + strconv.FormatInt(time.Now().UnixMilli(), 10)

	err = dataStore.StoreCards(playerID, cards)
	assert.NoError(t, err, "failed to store cards")

	storedCards, err := dataStore.GetCards(playerID)
	assert.NoError(t, err, "failed to get cards")

	assert.Len(t, storedCards, 1, "expected 1 card")
}

func TestInsertPlayer(t *testing.T) {
	dataStore := New("localhost", 5432, "postgres", "postgres", "progression")
	err := dataStore.Connect()
	assert.NoError(t, err)

	playerID := "test_player" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	player := Player{
		Id:        playerID,
		WildCards: 0,
		WildPacks: 0,
	}

	err = dataStore.UpdatePlayer(player)
	assert.NoError(t, err, "failed to store player")
}

func TestGetPlayer(t *testing.T) {
	dataStore := New("localhost", 5432, "postgres", "postgres", "progression")
	err := dataStore.Connect()
	assert.NoError(t, err)

	playerID := "test_player" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	player := Player{
		Id:        playerID,
		WildCards: 1,
		WildPacks: 23,
	}

	err = dataStore.UpdatePlayer(player)
	assert.NoError(t, err, "failed to store player")

	storedPlayer, err := dataStore.GetPlayer(playerID)
	assert.NoError(t, err, "failed to get player")
	assert.Equal(t, player, storedPlayer, "player did not match")
}

func TestUpdatePlayer(t *testing.T) {
	dataStore := New("localhost", 5432, "postgres", "postgres", "progression")
	err := dataStore.Connect()
	assert.NoError(t, err)

	playerID := "test_player" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	player := Player{
		Id:        playerID,
		WildCards: 1,
		WildPacks: 23,
	}

	err = dataStore.UpdatePlayer(player)
	assert.NoError(t, err, "failed to store player")

	player.WildCards = 0
	err = dataStore.UpdatePlayer(player)
	assert.NoError(t, err, "failed to store player")

	storedPlayer, err := dataStore.GetPlayer(playerID)
	assert.NoError(t, err, "failed to get player")
	assert.Equal(t, player, storedPlayer, "player did not match")
}
