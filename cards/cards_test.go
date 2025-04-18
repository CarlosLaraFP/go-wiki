package cards

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeck(t *testing.T) {
	deck := NewDeck(false)
	assert.Equal(t, 52, deck.Size())
	bottom := Card{Number: 1, Type: types[0]}
	assert.Equal(t, bottom, deck.Cards.Elements[0])
	top := Card{Number: 13, Type: types[len(types)-1]}
	assert.Equal(t, top, deck.Cards.Elements[len(deck.Cards.Elements)-1])
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
		//deck.Show()
	}
}

func TestSaveAndLoad(t *testing.T) {
	deck := NewDeck(false)

	if err := deck.Save(FilePath); err != nil {
		t.Errorf("failed to Save Deck: %v", err)
	}

	if deck, err := Load(FilePath); err != nil {
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
	//deck.Show()
	if deck.Size() != 52 {
		t.Errorf("expected %d cards; shuffle result %d", 52, deck.Size())
	}
}

func TestThreadSafety(t *testing.T) {
	deck := NewDeck(false)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		deck.Discard()
	}()
	go func() {
		defer wg.Done()
		deck.Deal(5)
	}()
	wg.Wait()
	assert.Equal(t, 47, deck.Size())

	newTop := Card{Number: 8, Type: types[len(types)-1]}
	top := Card{Number: 13, Type: types[len(types)-1]}

	if deck.Cards.Elements[0] == top {
		fmt.Println("Discard executed first")
	} else if deck.Cards.Elements[0] == newTop {
		fmt.Println("Deal executed first")
	} else {
		t.Errorf("Something went wrong.")
		deck.Show()
	}
}
