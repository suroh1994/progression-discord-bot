package repository

// DataStore is a backend for persisting the cards generated for every player.
type DataStore interface {
	// Connect connects the datastore to its respective backend. This doesn't necessarily entail any actions, but has to be called before the datastore can be used.
	Connect() error
	StartLeague() error
	EndLeague() error
	GetRound() (int, error)
	GetCards(userID string) ([]Card, error)
	StoreCards(userID string, cards []Card) error
	GetAllPlayers() ([]Player, error)
	GetPlayer(userID string) (Player, error)
	UpdatePlayer(player Player) error
	DropPlayer(userID string) error
	GetPairing(userID string) (Pairing, error)
	StorePairings(pairings []Pairing) error
	UpdatePairing(pairing Pairing) error
	IsAdmin(userID string) (bool, error)
	MakeAdmin(userID string) error
	GetBannedCards() ([]Ban, error)
	BanCard(cardName string) error
	UnbanCard(cardName string) error
}
