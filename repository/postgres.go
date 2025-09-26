package repository

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDataStore struct {
	hostname string
	port     int
	username string
	password string
	database string
	db       *gorm.DB
}

func New(hostname string, port int, username string, password string, database string) DataStore {
	return &postgresDataStore{
		hostname: hostname,
		port:     port,
		username: username,
		password: password,
		database: database,
	}
}

func (p *postgresDataStore) Connect() error {
	const errMsg = "unable to connect to datastore: %w"
	db, err := gorm.Open(postgres.Open(p.generateDSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	p.db = db
	return nil
}

func (p *postgresDataStore) generateDSN() string {
	// TODO Properly extract the current timezone of the server...
	TZ := "Europe/Berlin"
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		p.hostname, p.username, p.password, p.database, p.port, TZ)
}

func (p *postgresDataStore) StoreCards(userID string, cards []Card) error {
	const errMsg = "failed to store cards: %w"
	const query = `
			INSERT INTO player_card_pool (id, set_code, collector_number, count) VALUES %s
			ON CONFLICT (id, set_code, collector_number)
			DO UPDATE SET count = EXCLUDED.count + player_card_pool.count`

	fields, args := generateRows(userID, cards)
	result := p.db.Exec(fmt.Sprintf(query, fields), args...)

	if result.Error != nil {
		return fmt.Errorf(errMsg, result.Error)
	}

	return nil
}

func generateRows(userID string, cards []Card) (string, []any) {
	cardCounts := map[string]int{}

	// group similar cards
	for _, card := range cards {
		key := card.Set + "|" + strconv.Itoa(card.CollectorNumber)
		count, exists := cardCounts[key]
		if !exists {
			count = 0
		}

		cardCounts[key] = count + 1
	}

	// generate row per card
	inClause := make([]string, 0, len(cardCounts))
	args := make([]any, 0, len(cardCounts)*4)
	for key, count := range cardCounts {
		inClause = append(inClause, "(?, ?, ?, ?)")
		keyParts := strings.Split(key, "|")
		args = append(args, userID, keyParts[0], keyParts[1], count)
	}

	inClauseString := strings.Join(inClause, ", ")
	return inClauseString, args
}

func (p *postgresDataStore) GetCards(userID string) ([]Card, error) {
	const errMsg = "failed to fetch cards: %w"

	var cards []Card
	result := p.db.Table("player_card_pool").
		Where("id = ?", userID).
		Find(&cards)
	if result.Error != nil {
		return nil, fmt.Errorf(errMsg, result.Error)
	}

	return cards, nil
}

func (p *postgresDataStore) GetPlayer(userID string) (Player, error) {
	const errMsg = "failed to get player: %w"

	var player Player
	result := p.db.Table("player").First(&player, "id = ?", userID)
	if result.Error != nil {
		return player, fmt.Errorf(errMsg, result.Error)
	}

	return player, nil
}

func (p *postgresDataStore) UpdatePlayer(player Player) error {
	const errMsg = "failed to update player: %w"

	result := p.db.Table("player").Save(&player)
	if result.Error != nil {
		return fmt.Errorf(errMsg, result.Error)
	}

	return nil
}

func (p *postgresDataStore) StorePairings(pairings []Pairing) error {
	//TODO implement me
	panic("implement me")
}

func (p *postgresDataStore) UpdatePairing(userID string, wins, loses, draws int) error {
	//TODO implement me
	panic("implement me")
}
