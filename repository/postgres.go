package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"progression/repository/postgres"
)

type postgresDataStore struct {
	hostname string
	port     int
	username string
	password string
	database string
	db       *postgres.Queries
	conn     *pgx.Conn
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
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", p.username, p.password, p.hostname, p.port, p.database)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}
	p.conn = conn
	p.db = postgres.New(conn)
	return nil
}

func (p *postgresDataStore) StartLeague() error {
	const errMsg = "failed to start league: %w"

	_, err := p.GetRound()
	if err != nil && !errors.Is(err, ErrNoActiveLeague) {
		return fmt.Errorf(errMsg, err)
	}

	if err == nil {
		return fmt.Errorf(errMsg, ErrLeagueAlreadyOngoing)
	}

	err = p.db.StartLeague(context.Background())
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}
	return nil
}

func (p *postgresDataStore) EndLeague() error {
	const errMsg = "failed to end league: %w"
	err := p.db.EndLeague(context.Background())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "20000" { // not_found in custom exception
				return fmt.Errorf(errMsg, ErrNoActiveLeague)
			}
		}
		return fmt.Errorf(errMsg, err)
	}
	return nil
}

func (p *postgresDataStore) GetRound() (int, error) {
	const errMsg = "failed to get current round: %w"
	round, err := p.db.GetCurrentRound(context.Background())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf(errMsg, ErrNoActiveLeague)
		}
		return 0, fmt.Errorf(errMsg, err)
	}
	return int(round), nil
}

func (p *postgresDataStore) GetCards(userID string) ([]Card, error) {
	const errMsg = "failed to fetch cards: %w"
	cards, err := p.db.GetCards(context.Background(), userID)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	ret := make([]Card, len(cards))
	for i, c := range cards {
		ret[i] = Card{
			Name:            c.Name,
			Set:             c.SetCode,
			CollectorNumber: int(c.CollectorNumber),
			Count:           int(c.Count),
		}
	}
	return ret, nil
}

func (p *postgresDataStore) StoreCards(userID string, cards []Card) error {
	const errMsg = "failed to store cards: %w"

	tx, err := p.conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}
	defer tx.Rollback(context.Background())
	qtx := p.db.WithTx(tx)

	playerIDs := make([]string, len(cards))
	names := make([]string, len(cards))
	setCodes := make([]string, len(cards))
	collectorNumbers := make([]int32, len(cards))
	counts := make([]int32, len(cards))

	for i, card := range cards {
		playerIDs[i] = userID
		names[i] = card.Name
		setCodes[i] = card.Set
		collectorNumbers[i] = int32(card.CollectorNumber)
		counts[i] = int32(card.Count)
	}

	err = qtx.StoreCards(context.Background(), postgres.StoreCardsParams{
		PlayerID:        playerIDs,
		Name:            names,
		SetCode:         setCodes,
		CollectorNumber: collectorNumbers,
		Count:           counts,
	})

	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return tx.Commit(context.Background())
}

func (p *postgresDataStore) GetAllPlayers() ([]Player, error) {
	const errMsg = "failed to get players: %w"
	players, err := p.db.GetAllPlayers(context.Background())
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	ret := make([]Player, len(players))
	for i, p := range players {
		ret[i] = Player{
			Id:        p.ID,
			WildCards: int(p.WildCardCount),
			WildPacks: int(p.WildPackCount),
		}
	}
	return ret, nil
}

func (p *postgresDataStore) GetPlayer(userID string) (Player, error) {
	const errMsg = "failed to get player: %w"
	player, err := p.db.GetPlayer(context.Background(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Player{}, ErrPlayerNotFound
		}
		return Player{}, fmt.Errorf(errMsg, err)
	}
	return Player{
		Id:        player.ID,
		WildCards: int(player.WildCardCount),
		WildPacks: int(player.WildPackCount),
	}, nil
}

func (p *postgresDataStore) UpdatePlayer(player Player) error {
	const errMsg = "failed to update player: %w"
	err := p.db.UpdatePlayer(context.Background(), postgres.UpdatePlayerParams{
		ID:            player.Id,
		WildCardCount: int32(player.WildCards),
		WildPackCount: int32(player.WildPacks),
	})
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}
	return nil
}

func (p *postgresDataStore) GetPairing(userID string) (Pairing, error) {
	const errMsg = "failed to get pairing: %w"
	round, err := p.GetRound()
	if err != nil {
		return Pairing{}, fmt.Errorf(errMsg, err)
	}

	pairing, err := p.db.GetPairing(context.Background(), postgres.GetPairingParams{
		Round:     int32(round),
		PlayerId1: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Pairing{}, ErrPairingNotFound
		}
		return Pairing{}, fmt.Errorf(errMsg, err)
	}

	return Pairing{
		Round:     int(pairing.Round),
		PlayerId1: pairing.PlayerId1,
		PlayerId2: pairing.PlayerId2,
		Wins1:     int(pairing.Wins1),
		Wins2:     int(pairing.Wins2),
		Draws:     int(pairing.Draws),
	}, nil
}

func (p *postgresDataStore) StorePairings(pairings []Pairing) error {
	const errMsg = "failed to store pairings: %w"

	tx, err := p.conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}
	defer tx.Rollback(context.Background())
	qtx := p.db.WithTx(tx)

	rounds := make([]int32, len(pairings))
	player1IDs := make([]string, len(pairings))
	player2IDs := make([]string, len(pairings))
	wins1 := make([]int32, len(pairings))
	wins2 := make([]int32, len(pairings))
	draws := make([]int32, len(pairings))

	for i, p := range pairings {
		rounds[i] = int32(p.Round)
		player1IDs[i] = p.PlayerId1
		player2IDs[i] = p.PlayerId2
		wins1[i] = int32(p.Wins1)
		wins2[i] = int32(p.Wins2)
		draws[i] = int32(p.Draws)
	}

	err = qtx.StorePairings(context.Background(), postgres.StorePairingsParams{
		Round:     rounds,
		PlayerId1: player1IDs,
		PlayerId2: player2IDs,
		Wins1:     wins1,
		Wins2:     wins2,
		Draws:     draws,
	})

	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return tx.Commit(context.Background())
}

func (p *postgresDataStore) UpdatePairing(pairing Pairing) error {
	const errMsg = "failed to update pairing: %w"
	err := p.db.UpdatePairing(context.Background(), postgres.UpdatePairingParams{
		Round:     int32(pairing.Round),
		PlayerId1: pairing.PlayerId1,
		PlayerId2: pairing.PlayerId2,
		Wins1:     int32(pairing.Wins1),
		Wins2:     int32(pairing.Wins2),
		Draws:     int32(pairing.Draws),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "20000" {
				return fmt.Errorf(errMsg, ErrPairingNotFound)
			}
		}
		return fmt.Errorf(errMsg, err)
	}
	return nil
}