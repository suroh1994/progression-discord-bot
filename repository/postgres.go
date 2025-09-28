package repository

import (
	"errors"
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

func NewPostgresDataStore(hostname string, port int, username string, password string, database string) DataStore {
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
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Player{}, ErrPlayerNotFound
		}

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

func (p *postgresDataStore) GetPairing(userID string) (Pairing, error) {
	const errMsg = "failed to get pairing: %w"

	var pairing Pairing
	result := p.db.Table("pairing").
		Where("player1 = ?", userID).
		Or("player2 = ?", userID).
		Find(&pairing)
	if result.Error != nil {
		return pairing, fmt.Errorf(errMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return pairing, ErrPairingNotFound
	}

	return pairing, nil
}

func (p *postgresDataStore) StorePairings(pairings []Pairing) error {
	const errMsg = "failed to store pairings: %w"

	result := p.db.Table("pairing").Create(&pairings)
	if result.Error != nil {
		return fmt.Errorf(errMsg, result.Error)
	}

	return nil
}

func (p *postgresDataStore) UpdatePairing(pairing Pairing) error {
	const errMsg = "failed to update pairing: %w"

	const query = `UPDATE pairing SET wins1 = ?, wins2 = ?, draws = ?
               WHERE round = ? AND player1 = ? AND player2 = ?
               AND wins1 = 0 AND wins2 = 0 AND draws = 0`

	result := p.db.Exec(query, pairing.Wins1, pairing.Wins2, pairing.Draws, pairing.Round, pairing.Player1, pairing.Player2)
	if result.Error != nil {
		return fmt.Errorf(errMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf(errMsg, ErrPairingNotFound)
	}

	return nil
}

func (p *postgresDataStore) StartLeague() error {
	const errMsg = "failed to start league: %w"
	const query = `INSERT INTO league VALUES (0, true, now());`

	_, err := p.GetRound()
	if err != nil && !errors.Is(err, ErrNoActiveLeague) {
		return fmt.Errorf(errMsg, err)
	}

	if err == nil {
		return fmt.Errorf(errMsg, ErrLeagueAlreadyOngoing)
	}

	result := p.db.Exec(query)
	if result.Error != nil {
		return fmt.Errorf(errMsg, result.Error)
	}

	return nil
}

func (p *postgresDataStore) EndLeague() error {
	const errMsg = "failed to end league: %w"
	const query = `UPDATE league SET active = false WHERE active = true;`

	_, err := p.GetRound()
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	result := p.db.Exec(query)
	if result.Error != nil {
		return fmt.Errorf(errMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf(errMsg, ErrNoActiveLeague)
	}

	return nil
}

func (p *postgresDataStore) GetRound() (int, error) {
	const errMsg = "failed to get current round: %w"
	const query = `SELECT round FROM league where active = true;`

	var round int
	result := p.db.Raw(query).Find(&round)
	if result.Error != nil {
		return 0, fmt.Errorf(errMsg, result.Error)
	}

	if result.RowsAffected == 0 {
		return 0, fmt.Errorf(errMsg, ErrNoActiveLeague)
	}

	return round, nil
}
