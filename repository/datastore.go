package repository

// DataStore is a backend for persisting the cards generated for every player.
type DataStore interface {
	Connect() error
	StoreCards(userID string, cards []Card) error
	GetCards(userID string) ([]Card, error)
}
