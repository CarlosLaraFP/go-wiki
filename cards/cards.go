package cards

import "fmt"

type Card struct {
	Number  int
	Type    string
	IsJoker bool
}

var types = [4]string{"hearts", "spades", "diamonds", "tulips"}

// NewDeck returns a new deck of cards
func NewDeck(log bool) []Card {
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
				fmt.Printf("%v of %s | joker: %v\n", deck[i].Number, deck[i].Type, deck[i].IsJoker)
			}
		}
	}

	return deck
}
