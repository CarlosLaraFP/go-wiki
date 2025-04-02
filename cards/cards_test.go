package cards

import (
	"testing"
)

func TestNewDeck(t *testing.T) {
	deck := NewDeck(true)
	if len(deck.Cards) != 40 {
		t.Errorf("Deck should contain 40 cards")
	}

	card := deck.Deal()
	if card.Number != 9 && card.Type != "Hearts" && card.IsJoker {
		t.Errorf("Deal() did not pop the stack correctly: %v\n", card)
	}

	if len(deck.Cards) != 39 {
		t.Errorf("Deck should contain 39 cards")
	}
}
