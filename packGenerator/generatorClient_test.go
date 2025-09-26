package packGenerator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// IMPORTANT: (re-)start the magic booster pack generator with `make run-mbpg` before you run these tests

var client = New("http://localhost:8080")

func TestClient_CheckCard_exists(t *testing.T) {
	exists, err := client.CheckCard("IKO", 1)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestClient_CheckCard_does_not_exists(t *testing.T) {
	exists, err := client.CheckCard("IKO", 0)
	assert.NoError(t, err)
	assert.False(t, exists)

	exists, err = client.CheckCard("abcd", 1)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestClient_GetPacks_exists(t *testing.T) {
	cards, err := client.GetPacks("IKO", 1)
	assert.NoError(t, err)
	assert.Equal(t, 15, len(cards))
}

func TestClient_GetPacks_exists_multiple_packs(t *testing.T) {
	cards, err := client.GetPacks("IKO", 10)
	assert.NoError(t, err)
	assert.Equal(t, 150, len(cards))
}

// Apparently this is not yet implemented in the generator and therefore fails.
// Looking at the issues in the repository, there are even more rules, which haven't been implemented.
func TestClient_GetPacks_modern_set_with_less_cards(t *testing.T) {
	cards, err := client.GetPacks("FIN", 1)
	assert.NoError(t, err)
	assert.Equal(t, 14, len(cards))
}

func TestClient_GetPacks_does_not_exist(t *testing.T) {
	_, err := client.GetPacks("abcd", 1)
	assert.ErrorIs(t, err, ErrSetNotFound)
}
