package main

import (
	"reflect"
	"testing"
)

func TestRankSuitString(t *testing.T) {
	testCases := []struct {
		c Card
		r string
		s string
	}{
		{Card{Ace, Hearts}, "Ace", "Hearts"},
		{Card{Seven, Spades}, "Seven", "Spades"},
		{Card{Jack, Clubs}, "Jack", "Clubs"},
	}

	for _, tc := range testCases {
		card := tc.c
		if card.Suit.String() != tc.s || card.Rank.String() != tc.r {
			t.Errorf(`
generated string vals for Card %#v incorrect.
got Suit %v, Rank %v
want Suit %v, Rank %v
			`,
				card,
				card.Suit.String(), card.Rank.String(), tc.s, tc.r)
		}
	}
}

func TestNew(t *testing.T) {
	got := New()

	// The new Deck should have 52 cards.
	if len(got) != 52 {
		t.Errorf("length of deck should be 52, is %v", len(got))
	}

	// The new Deck should be sorted by Suit.
	want := New(SortBySuit)
	if !reflect.DeepEqual(got, want) {
		t.Errorf(`
		new deck not properly ordered by Suit and Rank.
		got %v
		want %v
		`, got, want)
	}
}

func TestAbsRank(t *testing.T) {
	testCases := []struct {
		c    Card
		want int
	}{
		{Card{Ace, Spades}, 1},
		{Card{Three, Clubs}, 29},
		{Card{Jack, Hearts}, 50},
	}

	for _, tc := range testCases {
		got := absRank(tc.c)
		if got != tc.want {
			t.Errorf("want %d got %d for the absolute value of card %#v", tc.want, got, tc.c)
		}
	}
}

func TestFunctionalOptions(t *testing.T) {
	t.Run("Shuffle to randomize order", func(t *testing.T) {
		// Chance that Shuffle Shuffles careds into exact same
		// order is possible but extremely unlikely.
		notWant := New()
		got := New(Shuffle)

		if reflect.DeepEqual(got, notWant) {
			t.Errorf(`
			unShuffled deck matches Shuffled deck.
			unShuffled deck %v
			Shuffled deck %v
			`, notWant, got)
		}
	})

	t.Run("sort by Rank", func(t *testing.T) {
		// Manually create a deck that's sorted by Rank.
		var want Deck
		suitID, rankID := 0, 0
		for rankID < 13 { // For each Rank...
			suitID = 0
			for suitID < 4 { // For each Suit...

				// Create a card and add it to the deck.
				c := Card{Rank(rankID), Suit(suitID)}
				want = append(want, c)
				suitID++
			}
			rankID++
		}

		got := New(SortByRank)

		if !reflect.DeepEqual(got, want) {
			t.Errorf(`
			deck not correctly sorted by Rank.
			got %v
			want %v
			`, got, want)
		}
	})

	t.Run("sort by Suit", func(t *testing.T) {
		// Decks are sorted by Rank by default.
		want := New()
		got := New(SortBySuit)

		if !reflect.DeepEqual(got, want) {
			t.Errorf(`
			deck not correctly sorted by Rank. 
			got %v
			want %v
			`, got, want)
		}
	})

	t.Run("add jokers test", func(t *testing.T) {
		got := New(AddJokers(4))
		want := New()

		// Manually add Four jokers to the 'want' deck.
		for i := 0; i < 4; i++ {
			want = append(want, Card{Rank(14), Suit(-1)})
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf(`
			Four jokers not correctly added to deck. 
			got %v
			wnat %v
			`, got, want)
		}

	})

	t.Run("filter out test", func(t *testing.T) {
		FilterOutJokers := FilterOut("Joker")
		got := New(AddJokers(4), FilterOutJokers)
		want := New()

		// Deck should be back to its original state after
		// adding and then removing jokers.
		if !reflect.DeepEqual(got, want) {
			t.Errorf(`
			jokers not successfully filtered from the deck.
			got %v
			want %v
			`, got, want)
		}
	})

	t.Run("add decks test", func(t *testing.T) {
		// Creat a function to add Two decks and
		// pass into NewDeck.
		addTwoDecks := AddDecks(2)
		got := New(addTwoDecks)

		// Test length 156.
		if len(got) != (52 * 3) {
			t.Errorf("length of deck with %v added decks is %v, should be %v", 2, len(got), 52*(3))
		}

		// Test against Three manually combined decks.
		want := New()
		want = append(want, New()...)
		want = append(want, New()...)

		if !reflect.DeepEqual(got, want) {
			t.Errorf(`
			%v decks not correctly added to new deck.
			got %v
			want %v
			`, 2, got, want)
		}

	})
}
