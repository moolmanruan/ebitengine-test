package deck

import (
	"math/rand"
	"time"
)

type Deck[T any] struct {
	cards []T
	r     *rand.Rand
}

// New creates a new deck consisting of the cards passed in.
func New[T any](cards []T) Deck[T] {
	cc := make([]T, len(cards))
	copy(cc, cards)
	return Deck[T]{cards: cc,
		r: rand.New(rand.NewSource(int64(time.Now().Nanosecond())))}
}

// Size returns the number of cards currently in the deck.
func (d Deck[T]) Size() int {
	return len(d.cards)
}

// Card returns the card at index i.
func (d Deck[T]) Card(i int) T {
	if len(d.cards) < i {
		var zero T
		return zero
	}
	return d.cards[i]
}

// Cards returns a copy of the list of cards in the deck.
func (d Deck[T]) Cards() []T {
	cc := make([]T, len(d.cards))
	copy(cc, d.cards)
	return cc
}

// Shuffle returns shuffled copy of the deck.
func (d Deck[T]) Shuffle() Deck[T] {
	cc := d.Cards()
	d.r.Shuffle(len(cc), func(i, j int) {
		cc[i], cc[j] = cc[j], cc[i]
	})
	return Deck[T]{cards: cc}
}

// Draw removes the top `num` cards from the deck and returns them along with
// a copy of the deck after the cards are removed.
func (d Deck[T]) Draw(num int) ([]T, Deck[T]) {
	if num > len(d.cards) {
		num = len(d.cards)
	}
	drawn := make([]T, num)
	copy(drawn, d.cards[:num])
	rest := make([]T, len(d.cards)-num)
	copy(rest, d.cards[num:])
	return drawn, Deck[T]{cards: rest}
}
