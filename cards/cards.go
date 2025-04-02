package cards

import (
	"fmt"
)

type Stack[T any] interface {
	Push(v T)
	Pop() *T
}

type Card struct {
	Number  int
	Type    string
	IsJoker bool
}

type Deck struct {
	Cards []Card
}

var types = [4]string{"Hearts", "Spades", "Diamonds", "Tulips"}

// NewDeck returns a new deck of cards
func NewDeck(log bool) Deck {
	deck := make([]Card, 0)
	isJoker := false

	for _, t := range types {
		for i := range 10 {
			if i == 0 {
				isJoker = true
			} else {
				isJoker = false
			}
			deck = append(deck, Card{
				Number:  i,
				Type:    t,
				IsJoker: isJoker,
			})
			if log {
				fmt.Printf("%v of %s | Joker: %v\n", deck[i].Number, deck[i].Type, deck[i].IsJoker)
			}
		}
	}
	return Deck{Cards: deck}
}

// Deal pops 1 element from the stop of the deck FIFO stack (O(1))
func (deck *Deck) Deal() *Card {
	if len(deck.Cards) == 0 {
		fmt.Println("Deck is empty. Game over.")
		return nil
	}
	top := &deck.Cards[len(deck.Cards)-1]
	deck.Cards = deck.Cards[:len(deck.Cards)-1]
	return top
}
