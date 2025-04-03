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
	s.elements = append([]T{element}, s.elements...)
}

func (s *Stack[T]) Size() int {
	return len(s.elements)
}

func (s *Stack[T]) Show() {
	for i := s.Size() - 1; i >= 0; i-- {
		fmt.Printf("%v\n", s.elements[i])
	}
}

type Card struct {
	Number  int
	Type    string
	IsJoker bool
}

// note type alias vs a new, distinct type (type Cards Stack[Card])
type Cards = Stack[Card]

// if Cards were a new type, only the methods defined specifically for the new type would be promoted to the parent struct
type Deck struct {
	Cards
}

// Hand returns the top n cards from the deck
func (d *Deck) Hand(n int) (*Cards, error) {
	var hand Cards

	for range n {
		card := d.Hit()
		if card == nil {
			return nil, fmt.Errorf("not enough cards left in deck to deal a full hand")
		}
		hand.PushTop(*card)
	}

	return &hand, nil
}

// Deal pops 1 element from the stop of the deck FIFO stack (O(1))
func (d *Deck) Hit() *Card {
	return d.Pop()
}

// Discard pops and pushes to the bottom
func (d *Deck) Discard() {
	if card := d.Pop(); card != nil {
		d.PushBottom(*card)
	}
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
