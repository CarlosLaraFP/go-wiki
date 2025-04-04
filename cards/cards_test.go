package cards

import (
	"fmt"
	"os"
	"testing"
)

func TestNewDeck(t *testing.T) {
	deck := NewDeck(false)
	if deck.Size() != 52 {
		t.Errorf("Deck should contain 52 cards")
	}

	card := deck.Hit()
	if card.Number != 13 && card.Type != "Clubs" {
		t.Errorf("Deal() did not pop the stack correctly: %v\n", card)
	}
	fmt.Printf("%v\n", card)
	if deck.Size() != 51 {
		t.Errorf("Deck should contain 51 cards")
	}

	deck.Discard()
	if deck.Size() != 51 {
		t.Errorf("Deck should still contain 51 cards")
	}
	deck.Show()

	hand, err := deck.Deal(5)
	if err != nil || (deck.Size() != 46 && len(hand) != 5) {
		t.Fail()
	}

	if err = deck.Save(FilePath); err != nil {
		t.Errorf("failed to Save Deck: %v", err)
	}

	if deck, err := LoadDeck(FilePath); err != nil {
		t.Errorf("failed to Load Deck: %v", err)

		if deck.Size() != 46 {
			t.Errorf("expected %d cards; loaded %d", 46, deck.Size())
		}
	}

	t.Cleanup(func() {
		os.Remove(FilePath)
	})
}
