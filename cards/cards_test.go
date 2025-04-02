package cards

import (
	"fmt"
	"testing"
)

func TestNewDeck(t *testing.T) {
	deck := NewDeck(true)
	fmt.Println(len(deck))
}
