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

func TestInsertPairing(t *testing.T) {
	dataStore := New("localhost", 5432, "postgres", "postgres", "progression")
	err := dataStore.Connect()
	assert.NoError(t, err)

	playerID := "test_player" + strconv.FormatInt(time.Now().UnixNano(), 10)
	playerID2 := "test_player" + strconv.FormatInt(time.Now().UnixNano(), 10)
	pairings := []Pairing{
		{
			Round:   1,
			Player1: playerID,
			Player2: playerID2,
			Wins1:   0,
			Wins2:   0,
			Draws:   0,
		},
	}

	err = dataStore.StorePairings(pairings)
	assert.NoError(t, err, "failed to store pairings")
}

func TestInsertPairings(t *testing.T) {
	dataStore := New("localhost", 5432, "postgres", "postgres", "progression")
	err := dataStore.Connect()
	assert.NoError(t, err)

	playerIDs := make([]string, 8)
	for i := range playerIDs {
		playerIDs[i] = "test_player" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	pairings := make([]Pairing, 0, len(playerIDs)/2)
	for i := 0; i < len(playerIDs)/2; i++ {
		pairings = append(pairings, Pairing{
			Round:   1,
			Player1: playerIDs[2*i],
			Player2: playerIDs[2*i+1],
			Wins1:   0,
			Wins2:   0,
			Draws:   0,
		})
	}

	err = dataStore.StorePairings(pairings)
	assert.NoError(t, err, "failed to store pairings")
}

func TestGetPairing_Player1(t *testing.T) {
	dataStore := New("localhost", 5432, "postgres", "postgres", "progression")
	err := dataStore.Connect()
	assert.NoError(t, err)

	playerIDs := make([]string, 8)
	for i := range playerIDs {
		playerIDs[i] = "test_player" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	pairings := make([]Pairing, 0, len(playerIDs)/2)
	for i := 0; i < len(playerIDs)/2; i++ {
		pairings = append(pairings, Pairing{
			Round:   1,
			Player1: playerIDs[2*i],
			Player2: playerIDs[2*i+1],
			Wins1:   0,
			Wins2:   0,
			Draws:   0,
		})
	}

	err = dataStore.StorePairings(pairings)
	assert.NoError(t, err, "failed to store pairings")

	pairing, err := dataStore.GetPairing(playerIDs[2])
	assert.NoError(t, err, "failed to get pairing")
	assert.Equal(t, pairing, pairings[1], "pairing did not match")
}

func TestGetPairing_Player2(t *testing.T) {
	dataStore := New("localhost", 5432, "postgres", "postgres", "progression")
	err := dataStore.Connect()
	assert.NoError(t, err)

	playerIDs := make([]string, 8)
	for i := range playerIDs {
		playerIDs[i] = "test_player" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	pairings := make([]Pairing, 0, len(playerIDs)/2)
	for i := 0; i < len(playerIDs)/2; i++ {
		pairings = append(pairings, Pairing{
			Round:   1,
			Player1: playerIDs[2*i],
			Player2: playerIDs[2*i+1],
			Wins1:   0,
			Wins2:   0,
			Draws:   0,
		})
	}

	err = dataStore.StorePairings(pairings)
	assert.NoError(t, err, "failed to store pairings")

	pairing, err := dataStore.GetPairing(playerIDs[5])
	assert.NoError(t, err, "failed to get pairing")
	assert.Equal(t, pairing, pairings[2], "pairing did not match")
}

func TestUpdatePairing(t *testing.T) {
	dataStore := New("localhost", 5432, "postgres", "postgres", "progression")
	err := dataStore.Connect()
	assert.NoError(t, err)

	playerID := "test_player" + strconv.FormatInt(time.Now().UnixNano(), 10)
	playerID2 := "test_player" + strconv.FormatInt(time.Now().UnixNano(), 10)
	pairings := []Pairing{
		{
			Round:   1,
			Player1: playerID,
			Player2: playerID2,
			Wins1:   0,
			Wins2:   0,
			Draws:   0,
		},
	}

	err = dataStore.StorePairings(pairings)
	assert.NoError(t, err, "failed to store pairings")

	pairings[0].Wins1 = 2
	err = dataStore.UpdatePairing(pairings[0])
	assert.NoError(t, err, "failed to update pairing")

	storedPairing, err := dataStore.GetPairing(playerID2)
	assert.NoError(t, err, "failed to get pairing")
	assert.Equal(t, pairings[0], storedPairing, "pairing did not match")
}
