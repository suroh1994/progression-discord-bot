package packGenerator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	hostAddress string
}

func New(hostAddress string) *Client {
	return &Client{
		hostAddress: hostAddress,
	}
}

func (c *Client) GetPacks(setCode string, count int) ([]Card, error) {
	const errMsg = "unable to get packs from generator: %w"

	url := fmt.Sprintf("%s/pack/%s?count=%d&export=false&outputformat=json&tokens=false",
		c.hostAddress, setCode, count)
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		fmt.Printf("failed to find set %q, got response code %d\n", setCode, response.StatusCode)
		return nil, fmt.Errorf(errMsg, ErrSetNotFound)
	}

	var cards []Card
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&cards)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return cards, nil
}

func (c *Client) CheckCard(setCode string, collectorNumber int) (bool, error) {
	const errMsg = "unable to get card from generator: %w"

	url := fmt.Sprintf("%s/card/%s/%d", c.hostAddress, setCode, collectorNumber)
	response, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf(errMsg, err)
	}

	return response.StatusCode == 200, nil
}
