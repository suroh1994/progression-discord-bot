package repository

// DataStore is a backend for persisting the cards generated for every player.
type DataStore interface {
	// Connect connects the datastore to its respective backend. This doesn't necessarily entail any actions, but has to be called before the datastore can be used.
	Connect() error
	GetCards(userID string) ([]Card, error)
	StoreCards(userID string, cards []Card) error
	GetPlayer(userID string) (Player, error)
	UpdatePlayer(player Player) error
	StorePairings(pairings []Pairing) error
	UpdatePairing(userID string, wins, loses, draws int) error
}
