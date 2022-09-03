package card

type Suit int

const (
	NoSuit   Suit = 0
	Clubs    Suit = 1
	Diamonds Suit = 2
	Hearts   Suit = 3
	Spades   Suit = 4
)

type Face int

const (
	Ace   Face = 1
	Two   Face = 2
	Three Face = 3
	Four  Face = 4
	Five  Face = 5
	Six   Face = 6
	Seven Face = 7
	Eight Face = 8
	Nine  Face = 9
	Ten   Face = 10
	Jack  Face = 11
	Queen Face = 12
	King  Face = 13
	Joker Face = 14
)

var Suits = []Suit{Clubs, Diamonds, Hearts, Spades}
var Faces = []Face{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

type Card struct {
	Suit Suit
	Face Face
}

// New creates a new Card struct with the supplied suit and face values
func New(suit Suit, face Face) Card {
	return Card{Suit: suit, Face: face}
}

// StandardDeck returns a list of cards you would find in a standard playing
// deck.
func StandardDeck() []Card {
	cc := make([]Card, 0, len(Suits)*len(Faces))
	for _, s := range Suits {
		for _, f := range Faces {
			cc = append(cc, New(s, f))
		}
	}
	return cc
}

// StandardDeckWithJokers returns the same as StandardDeck with an additional
// `n` joker cards.
func StandardDeckWithJokers(n int) []Card {
	cc := StandardDeck()
	for i := 0; i < n; i++ {
		cc = append(cc, New(NoSuit, Joker))
	}
	return cc
}
