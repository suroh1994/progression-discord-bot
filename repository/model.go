package repository

// Player represents a player in the league.
type Player struct {
	Id        string
	WildCards int `gorm:"column:wild_card_count"`
	WildPacks int `gorm:"column:wild_pack_count"`
	Dropped   bool
}

// Card represents a card in a players card pool.
type Card struct {
	Name            string
	Set             string
	CollectorNumber int
	Count           int
}

// Pairing represents a pairing of players in a round. Once any scores have been reported, the pairing is assumed to be over.
type Pairing struct {
	Round   int    `gorm:"primaryKey"`
	Player1 string `gorm:"primaryKey"`
	Player2 string `gorm:"primaryKey"`
	Wins1   int
	Wins2   int
	Draws   int
}
