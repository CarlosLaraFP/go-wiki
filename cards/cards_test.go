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
}

func TestHit(t *testing.T) {
	deck := NewDeck(false)
	card := deck.Hit()

	if card.Number != 13 && card.Type != "Clubs" {
		t.Errorf("Deal() did not pop the stack correctly: %v\n", card)
	}
	fmt.Printf("%v\n", card)
	if deck.Size() != 51 {
		t.Errorf("Deck should contain 51 cards")
	}
}

func TestDiscard(t *testing.T) {
	deck := NewDeck(false)
	deck.Discard()
	if deck.Size() != 52 {
		t.Errorf("Deck should still contain 52 cards")
	}
	deck.Show()
}

func TestDeal(t *testing.T) {
	deck := NewDeck(false)
	if hand, err := deck.Deal(5); err != nil || (deck.Size() != 47 && len(hand) != 5) {
		t.Errorf("Deal method failed")
		deck.Show()
	}
}

func TestSave_Load(t *testing.T) {
	deck := NewDeck(false)

	if err := deck.Save(FilePath); err != nil {
		t.Errorf("failed to Save Deck: %v", err)
	}

	if deck, err := LoadDeck(FilePath); err != nil {
		t.Errorf("failed to Load Deck: %v", err)

		if deck.Size() != 52 {
			t.Errorf("expected %d cards; loaded %d", 52, deck.Size())
		}
	}

	t.Cleanup(func() {
		os.Remove(FilePath)
	})
}

func TestShuffle(t *testing.T) {
	deck := NewDeck(false)
	deck.Show()
	deck.Shuffle()
	fmt.Println("Shuffled")
	deck.Show()
	if deck.Size() != 52 {
		t.Errorf("expected %d cards; shuffle result %d", 52, deck.Size())
	}
}
