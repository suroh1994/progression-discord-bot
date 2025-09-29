package repository

import (
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// IMPORTANT: (re-)start the database with `make run-pgdb` before you run these tests

var dataStore DataStore

func setupSuite() {
	dataStore = NewPostgresDataStore("localhost", 5432, "postgres", "postgres", "progression")
	err := dataStore.Connect()
	if err != nil {
		log.Fatal(err)
	}
}

func teardownSuite() {
}

func setupTest() {
	pgDS, ok := dataStore.(*postgresDataStore)
	if !ok {
		log.Fatal("setup test fail")
	}
	if err := pgDS.db.Exec("TRUNCATE TABLE league;").Error; err != nil {
		log.Fatal(err)
	}
	if err := pgDS.db.Exec("TRUNCATE TABLE player;").Error; err != nil {
		log.Fatal(err)
	}
	if err := pgDS.db.Exec("TRUNCATE TABLE player_card_pool;").Error; err != nil {
		log.Fatal(err)
	}
	if err := pgDS.db.Exec("TRUNCATE TABLE pairing;").Error; err != nil {
		log.Fatal(err)
	}
}

func teardownTest() {
}

func TestMain(m *testing.M) {
	setupSuite()
	code := m.Run()
	teardownSuite()
	os.Exit(code)
}

func TestInsertCardPool(t *testing.T) {
	setupTest()

	cards := []Card{
		{
			Name:            "Adaptive Shimmerer",
			Set:             "IKO",
			CollectorNumber: 1,
		},
	}
	playerID := "test_player" + strconv.FormatInt(time.Now().UnixMilli(), 10)

	err := dataStore.StoreCards(playerID, cards)
	assert.NoError(t, err, "failed to store cards")
	teardownTest()
}

func TestGetCardPool(t *testing.T) {
	setupTest()

	cards := []Card{
		{
			Name:            "Adaptive Shimmerer",
			Set:             "IKO",
			CollectorNumber: 1,
		},
		{
			Name:            "Farfinder",
			Set:             "IKO",
			CollectorNumber: 2,
		},
		{
			Name:            "Adaptive Shimmerer",
			Set:             "IKO",
			CollectorNumber: 1,
		},
	}
	playerID := "test_player" + strconv.FormatInt(time.Now().UnixMilli(), 10)

	err := dataStore.StoreCards(playerID, cards)
	assert.NoError(t, err, "failed to store cards")

	storedCards, err := dataStore.GetCards(playerID)
	assert.NoError(t, err, "failed to get cards")

	assert.Len(t, storedCards, 2, "expected 2 different cards")
	teardownTest()
}

func TestCardPoolDeduplicate(t *testing.T) {
	setupTest()

	cards := []Card{
		{
			Name:            "Adaptive Shimmerer",
			Set:             "IKO",
			CollectorNumber: 1,
		},
		{
			Name:            "Adaptive Shimmerer",
			Set:             "IKO",
			CollectorNumber: 1,
		},
	}
	playerID := "test_player" + strconv.FormatInt(time.Now().UnixMilli(), 10)

	err := dataStore.StoreCards(playerID, cards)
	assert.NoError(t, err, "failed to store cards")

	storedCards, err := dataStore.GetCards(playerID)
	assert.NoError(t, err, "failed to get cards")

	assert.Len(t, storedCards, 1, "expected 1 card")
	assert.Equal(t, 2, storedCards[0].Count, "expected 2 copies")
	teardownTest()
}

func TestInsertPlayer(t *testing.T) {
	setupTest()

	playerID := "test_player" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	player := Player{
		Id:        playerID,
		WildCards: 0,
		WildPacks: 0,
	}

	err := dataStore.UpdatePlayer(player)
	assert.NoError(t, err, "failed to store player")
	teardownTest()
}

func TestGetPlayer(t *testing.T) {
	setupTest()

	playerID := "test_player" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	player := Player{
		Id:        playerID,
		WildCards: 1,
		WildPacks: 23,
	}

	err := dataStore.UpdatePlayer(player)
	assert.NoError(t, err, "failed to store player")

	storedPlayer, err := dataStore.GetPlayer(playerID)
	assert.NoError(t, err, "failed to get player")
	assert.Equal(t, player, storedPlayer, "player did not match")
	teardownTest()
}

func TestUpdatePlayer(t *testing.T) {
	setupTest()

	playerID := "test_player" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	player := Player{
		Id:        playerID,
		WildCards: 1,
		WildPacks: 23,
	}

	err := dataStore.UpdatePlayer(player)
	assert.NoError(t, err, "failed to store player")

	player.WildCards = 0
	err = dataStore.UpdatePlayer(player)
	assert.NoError(t, err, "failed to store player")

	storedPlayer, err := dataStore.GetPlayer(playerID)
	assert.NoError(t, err, "failed to get player")
	assert.Equal(t, player, storedPlayer, "player did not match")
	teardownTest()
}

func TestInsertPairing(t *testing.T) {
	setupTest()

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

	err := dataStore.StorePairings(pairings)
	assert.NoError(t, err, "failed to store pairings")
	teardownTest()
}

func TestInsertPairings(t *testing.T) {
	setupTest()

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

	err := dataStore.StorePairings(pairings)
	assert.NoError(t, err, "failed to store pairings")
	teardownTest()
}

func TestGetPairing_Player1(t *testing.T) {
	setupTest()

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

	err := dataStore.StorePairings(pairings)
	assert.NoError(t, err, "failed to store pairings")

	pairing, err := dataStore.GetPairing(playerIDs[2])
	assert.NoError(t, err, "failed to get pairing")
	assert.Equal(t, pairing, pairings[1], "pairing did not match")
	teardownTest()
}

func TestGetPairing_Player2(t *testing.T) {
	setupTest()

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

	err := dataStore.StorePairings(pairings)
	assert.NoError(t, err, "failed to store pairings")

	pairing, err := dataStore.GetPairing(playerIDs[5])
	assert.NoError(t, err, "failed to get pairing")
	assert.Equal(t, pairing, pairings[2], "pairing did not match")
	teardownTest()
}

func TestUpdatePairing(t *testing.T) {
	setupTest()

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

	err := dataStore.StorePairings(pairings)
	assert.NoError(t, err, "failed to store pairings")

	pairings[0].Wins1 = 2
	err = dataStore.UpdatePairing(pairings[0])
	assert.NoError(t, err, "failed to update pairing")

	storedPairing, err := dataStore.GetPairing(playerID2)
	assert.NoError(t, err, "failed to get pairing")
	assert.Equal(t, pairings[0], storedPairing, "pairing did not match")
	teardownTest()
}

func TestStartRound(t *testing.T) {
	setupTest()

	err := dataStore.StartLeague()
	assert.NoError(t, err, "failed to start league")

	err = dataStore.StartLeague()
	assert.ErrorIs(t, err, ErrLeagueAlreadyOngoing, "failed to start league")
	teardownTest()
}

func TestEndRound(t *testing.T) {
	setupTest()

	err := dataStore.EndLeague()
	assert.ErrorIs(t, err, ErrNoActiveLeague, "ending league without an active league shouldn't work")

	err = dataStore.StartLeague()
	assert.NoError(t, err, "failed to start league")

	err = dataStore.EndLeague()
	assert.NoError(t, err, "failed to end league")

	err = dataStore.EndLeague()
	assert.ErrorIs(t, err, ErrNoActiveLeague, "ending league without an active league shouldn't work")
	teardownTest()
}
