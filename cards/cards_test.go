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
	bottom := Card{Number: 1, Type: types[0]}
	if deck.Cards.Elements[0] != bottom {
		t.Errorf("bottom card should be ace of hearts")
	}
	top := Card{Number: 13, Type: types[len(types)-1]}
	if deck.Cards.Elements[len(deck.Cards.Elements)-1] != top {
		t.Errorf("top card should be king of clubs")
	}
}

func TestHit(t *testing.T) {
	deck := NewDeck(false)
	top := Card{Number: 12, Type: types[len(types)-1]}
	card := deck.Hit()

	if card.Number != 13 && card.Type != "Clubs" {
		t.Errorf("Deal() did not pop the stack correctly: %v\n", card)
	}
	if deck.Size() != 51 {
		t.Errorf("Deck should contain 51 cards")
	}
	if deck.Cards.Elements[len(deck.Cards.Elements)-1] != top {
		t.Errorf("top card should now be queen of clubs")
	}
}

func TestDiscard(t *testing.T) {
	deck := NewDeck(false)
	top := Card{Number: 13, Type: types[len(types)-1]}
	deck.Discard()
	if deck.Size() != 52 {
		t.Errorf("Deck should still contain 52 cards")
	}
	if deck.Cards.Elements[0] != top {
		t.Errorf("bottom card should now be king of clubs")
	}
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
	deck.Shuffle()
	fmt.Println("Shuffled")
	deck.Show()
	if deck.Size() != 52 {
		t.Errorf("expected %d cards; shuffle result %d", 52, deck.Size())
	}
}
