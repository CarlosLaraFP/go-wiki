package cards

import (
	"fmt"
)

type Stack[T any] struct {
	elements []T
}

func (s *Stack[T]) Pop() *T {
	if len(s.elements) == 0 {
		return nil
	}
	top := &s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return top
}

func (s *Stack[T]) PushTop(element T) {
	s.elements = append(s.elements, element)
}

func (s *Stack[T]) PushBottom(element T) {
	deck := make([]T, 0)
	deck = append(deck, element)
	s.elements = append(deck, s.elements...)
}

func (s *Stack[T]) Size() int {
	return len(s.elements)
}

type Card struct {
	Number  int
	Type    string
	IsJoker bool
}

type Cards Stack[Card]

type Deck struct {
	Stack[Card]
}

// Deal pops 1 element from the stop of the deck FIFO stack (O(1))
func (deck *Deck) Deal() *Card {
	return deck.Pop()
}

// Discard pops and pushes to the bottom
func (deck *Deck) Discard() {
	if card := deck.Pop(); card != nil {
		deck.PushBottom(*card)
	}
}

func (deck *Deck) Size() int {
	return len(deck.elements)
}

var types = [4]string{"Hearts", "Spades", "Diamonds", "Clubs"}

// NewDeck returns a new deck of cards
func NewDeck(log bool) Deck {
	deck := make([]Card, 0)

	for _, t := range types {
		for i := range 13 {
			deck = append(deck, Card{
				Number: i + 1,
				Type:   t,
			})
			if log {
				fmt.Printf("%v of %s\n", deck[i].Number, deck[i].Type)
			}
		}
	}
	return Deck{Stack[Card]{elements: deck}}
}
