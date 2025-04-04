package cards

import (
	"encoding/json"
	"fmt"
	"os"
)

type Stack[T any] struct {
	Elements []T `json:"elements"`
}

func (s *Stack[T]) Pop() *T {
	if len(s.Elements) == 0 {
		return nil
	}
	top := &s.Elements[len(s.Elements)-1]
	s.Elements = s.Elements[:len(s.Elements)-1]
	return top
}

func (s *Stack[T]) PushTop(element T) {
	s.Elements = append(s.Elements, element)
}

func (s *Stack[T]) PushBottom(element T) {
	s.Elements = append([]T{element}, s.Elements...)
}

func (s *Stack[T]) Size() int {
	return len(s.Elements)
}

func (s *Stack[T]) Show() {
	for i := s.Size() - 1; i >= 0; i-- {
		fmt.Printf("%v\n", s.Elements[i])
	}
}

type Card struct {
	Number  int    `json:"number"`
	Type    string `json:"type"`
	IsJoker bool   `json:"isJoker,omitempty"`
}

type Hand []Card

// note type alias vs a new, distinct type
type Cards = Stack[Card]

// if Cards were a new type, only the methods defined specifically for the new type would be promoted to the parent struct
type Deck struct {
	Cards `json:"cards"`
}

// Deal returns the top n cards from the deck as a Hand
func (d *Deck) Deal(handSize int) (Hand, error) {
	var hand Hand

	for range handSize {
		card := d.Hit()
		if card == nil {
			return nil, fmt.Errorf("not enough cards left in deck to deal a full hand")
		}
		hand = append(hand, *card)
	}

	return hand, nil
}

// Deal pops 1 element from the stop of the deck FIFO stack (O(1) time complexity).
func (d *Deck) Hit() *Card {
	return d.Pop()
}

// Discard pops and pushes to the bottom (O(N) time complexity).
// Maintaining a secondary Deck would make Discard O(1) time complexity.
func (d *Deck) Discard() {
	if card := d.Pop(); card != nil {
		d.PushBottom(*card)
	}
}

func (d *Deck) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating OS file %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	if err := encoder.Encode(d); err != nil {
		return fmt.Errorf("error writing JSON to file %v", err)
	}
	return nil
}

func LoadDeck(path string) (*Deck, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %v", err)
	}
	defer file.Close()

	deck := Deck{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&deck); err != nil {
		return nil, fmt.Errorf("failed to read JSON from %v", err)
	}
	return &deck, nil
}

const FilePath = "deck.json"

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
	return Deck{Stack[Card]{Elements: deck}}
}
