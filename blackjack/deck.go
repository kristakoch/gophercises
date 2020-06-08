package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Suit represents the Suit value of the card.
type Suit int

// Suit constants.
const (
	Spades   Suit = iota // 0
	Diamonds             // 1
	Clubs                // 2
	Hearts               // 3
)

// Rank represents the number/face value of the card.
type Rank int

// Rank constants.
const (
	Ace Rank = iota + 1 // 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack  // 11
	Queen // 12
	King  // 13
	Joker // 14
)

// Card represents a playing card.
type Card struct {
	Rank
	Suit
}

func (c Card) String() string {
	if c.Rank.String() == "Joker" {
		return "Joker"
	}
	return fmt.Sprintf("%s of %s",
		c.Rank.String(),
		c.Suit.String())
}

// Deck represents a deck of playing cards.
type Deck []Card

// New is a factory function for Decks.
func New(options ...func(*Deck)) Deck {
	d := generateDeck()
	for _, option := range options {
		option(&d)
	}

	return d
}

// generateDeck does the work of populating the
// Deck with cards.
func generateDeck() Deck {
	var d Deck
	suitID, rankID := 0, 1
	for suitID < 4 { // For each Suit...
		rankID = 1
		for rankID < 14 { // For each Rank...

			// Create a card and add it to the deck.
			c := Card{Rank(rankID), Suit(suitID)}
			d = append(d, c)

			rankID++
		}
		suitID++
	}
	return d
}

// absRank gets the value of a Card
func absRank(c Card) int {
	return (13 * int(c.Suit)) + int(c.Rank)
}

// ===============================================
// FUNCTIONAL OPTIONS for use in NEWDECK
// ===============================================

// AddDecks is an option for adding a number of
//  Decks to the default Deck.
var AddDecks = func(extraDecks int) func(*Deck) {

	adder := func(deck *Deck) {
		ret := New()
		// Add the decks.
		for i := 0; i < extraDecks; i++ {
			ret = append(ret, New()...)
		}
		*deck = ret
	}

	return adder
}

// FilterOut filters out a Suit or Rank of cards.
var FilterOut = func(target string) func(*Deck) {

	filterer := func(deck *Deck) {
		var ret []Card
		for _, card := range *deck {
			r, s := card.Rank.String(), card.Suit.String()
			// Don't add cards with filter-out vals to new slice.
			if r == target || s == target {
				continue
			}
			ret = append(ret, card)
		}
		*deck = ret
	}

	return filterer
}

// Shuffle randomizes the Deck order.
var Shuffle = func(deck *Deck) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]Card, len(*deck))
	perm := r.Perm(len(*deck))
	for i, randIndex := range perm {
		ret[i] = (*deck)[randIndex]
	}
	*deck = ret
}

// SortBy takes in a user function to sort the deck by
// and returns a function to sort the deck by it.
var SortBy = func(lessFn func(i, j int) bool) func(*Deck) {
	sorter := func(deck *Deck) {
		sort.Slice(*deck, lessFn)
	}
	return sorter
}

// SortBySuit is an option to sort and group the
// cards by Suit.
var SortBySuit = func(deck *Deck) {
	sort.Slice(*deck, func(i, j int) bool {
		one, Two := (*deck)[i], (*deck)[j]
		if one.Suit == Two.Suit {
			return one.Rank < Two.Rank
		}
		return one.Suit < Two.Suit
	})
}

// SortByRank is an option to sort and group the
// cards by Rank.
var SortByRank = func(deck *Deck) {
	sort.Slice(*deck, func(i, j int) bool {
		one, Two := (*deck)[i], (*deck)[j]
		if one.Rank == Two.Rank {
			return one.Suit < Two.Suit
		}
		return one.Rank < Two.Rank
	})
}

// AddJokers adds 0-4 jokers to a Deck.9
var AddJokers = func(numJokers int) func(*Deck) {
	jokerAdder := func(deck *Deck) {
		for i := 0; i < numJokers; i++ {
			Joker := Card{Rank(14), Suit(-1)}
			*deck = append(*deck, Joker)
		}
	}
	return jokerAdder
}

// ---------------------out of use--------------------- //

// String implements stringer for the Suit type.
// note: commented out after use of stringer -type=Suit
// func (s Suit) String() string {
// 	switch s {
// 	case 0:
// 		return "Spades"
// 	case 1:
// 		return "Diamonds"
// 	case 2:
// 		return "Clubs"
// 	case 3:
// 		return "Hearts"
// 	default:
// 		return "Suit n/a"
// 	}
// }

// String implements stringer for the Rank type.
// note: commented out after use of stringer -type=Rank
// func (r Rank) String() string {
// 	ranks := []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King", "Joker"}
// 	return ranks[r]
// }
