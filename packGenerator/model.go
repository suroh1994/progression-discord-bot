package packGenerator

type Card struct {
	Name            string `json:"name"`
	Foil            bool   `json:"foil"`
	ScryfallURI     string `json:"scryfallURI"`
	Set             string `json:"set"`
	CollectorNumber string `json:"collectorNumber"`
	ImageURL        string `json:"imageURL"`
}
