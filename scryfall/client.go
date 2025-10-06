package scryfall

import (
	"context"
	"fmt"
	"strings"

	sf "github.com/BlueMonday/go-scryfall"
)

type Client struct {
	client         *sf.Client
	defaultOptions sf.SearchCardsOptions
}

type SearchOptionsModifier func(sf.SearchCardsOptions) sf.SearchCardsOptions

func NewClient() (*Client, error) {
	client, err := sf.NewClient()
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
		defaultOptions: sf.SearchCardsOptions{
			Unique:        sf.UniqueModeCards,
			Order:         sf.OrderSet,
			Dir:           sf.DirDesc,
			IncludeExtras: true,
		},
	}, nil
}

func (c *Client) SearchCardInSets(ctx context.Context, cardName string, sets []string, optionsModifiers ...SearchOptionsModifier) ([]sf.Card, error) {
	options := c.defaultOptions
	for _, modifier := range optionsModifiers {
		options = modifier(options)
	}

	setRestriction := generateSetRestriction(sets)

	query := fmt.Sprintf("%s (%s)", cardName, setRestriction)
	result, err := c.client.SearchCards(ctx, query, options)
	if err != nil {
		return nil, err
	}

	return result.Cards, nil
}

func (c *Client) SearchCard(ctx context.Context, cardName string, optionsModifiers ...SearchOptionsModifier) ([]sf.Card, error) {
	return c.SearchCardInSets(ctx, cardName, nil, optionsModifiers...)
}

// generateSetRestriction generates a string containing all given sets, which can be appended to a search to limit the results to only valid sets.
// TODO: the search query has a length limitation of 1000 charaters. To prevent errors for long leagues, we should return multiple strings if needed.
func generateSetRestriction(sets []string) string {
	if len(sets) == 0 {
		return ""
	}

	setSearchStrings := make([]string, len(sets))
	for i, set := range sets {
		setSearchStrings[i] = fmt.Sprintf("s:%s", set)
	}
	setRestriction := strings.Join(setSearchStrings, " or ")
	return "(" + setRestriction + ")"
}
