package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	// Create the start state for the game.
	g := NewGame(2)

	// Deliver the blackjack game.
	g.Play()
}

// player represents a non-dealer game member.
type player struct {
	name string
	hand []Card
}

// Game represents a game of blackjack.
type Game struct {
	gameDeck Deck
	dealer   struct {
		hand   []Card
		hidden Card
	}
	numPlayers    int
	players       []player
	turn          int
	playersBusted int
	playersWon    int
}

// NewGame creates a new instance of the game.
func NewGame(numPlayers int) *Game {
	var g Game
	// Create and set the game
	g.gameDeck = New(Shuffle)
	g.numPlayers = numPlayers

	// Deal out the cards.
	g.initialDeal()

	return &g
}

// Play initiates the blackjack game.
func (g *Game) Play() {
	// Check for a dealer blackjack.
	dealerHand := append(g.dealer.hand, g.dealer.hidden)
	if hasBlackjack(dealerHand) {
		fmt.Println("~ * dealer has blackjack! * ~")
		time.Sleep(time.Second * 1)

		// Game over.
		g.printFinalScores()
		return
	}

	// Check for player blackjacks.
	for _, p := range g.players {
		if handTotal(p.hand) == 21 {
			fmt.Println("~ * player", p.name, "has blackjack! * ~")
			time.Sleep(time.Second * 1)

			g.playersWon++
		}
	}

	// Start on the first turn and go to 2.
	for turn := 1; turn < 3; turn++ {
		g.turn = turn

		fmt.Println("========================")
		fmt.Println("TURN", g.turn)
		fmt.Println("========================")

		// Give each player a turn.
		for i, p := range g.players {
			total := handTotal(g.players[i].hand)
			if total < 21 {
				g.players[i].hand = g.takeTurn(p.hand, p.name)

				total = handTotal(g.players[i].hand)
				if total == 21 {
					// Player wins.
					fmt.Printf("~ * player %v wins! * ~\n", p.name)
					time.Sleep(time.Second * 1)
					g.playersWon++
				}
				if total > 21 {
					// Player busts.
					fmt.Printf("x x player %v busts! x x\n", p.name)
					time.Sleep(time.Second * 1)
					g.playersBusted++
				}
			}
		}

		// At the end of the first turn, reveal hidden dealer card.
		if g.turn == 1 {
			g.dealer.hand = append(g.dealer.hand, g.dealer.hidden)
			fmt.Printf("dealer reveals hidden card: %v\n", g.dealer.hidden)
			time.Sleep(time.Second * 1)
		}

		// At the end of the second turn, dealer
		// gets to have a turn.
		if g.turn == 2 {
			if (g.playersBusted + g.playersWon) != g.numPlayers {
				g.dealer.hand = g.takeTurn(g.dealer.hand, "dealer")
				total := handTotal(g.dealer.hand)
				if total == 21 {
					// Dealer wins..
					fmt.Println("~ * dealer wins! * ~")
					time.Sleep(time.Second * 1)
				}
				if total > 21 {
					// Dealer busts.
					fmt.Println("x x dealer busts! x x")
					time.Sleep(time.Second * 1)
				}
			} else {
				time.Sleep(time.Second * 1)
			}
		}
	}

	// Final scores...
	g.printFinalScores()
}

// takeTurn returns the final hand after player
// turn decisions to hit or stand.
func (g *Game) takeTurn(hand []Card, pName string) []Card {
	g.printPlayerHands()
	for {
		fmt.Printf(`
	%v,
	your hand total is %v
	enter h for hit or s for stand
	`, pName, handTotal(hand))

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		answer := scanner.Text()

		if answer == "h" {
			hand = append(
				hand,
				oneOffTheTop(&g.gameDeck),
			)

			fmt.Println("\nhand after hitting is...")
			for _, c := range hand {
				fmt.Println(c)
			}

			// Check for win or bust.
			// If so, break.
			total := handTotal(hand)
			if total > 20 {
				break
			}
		}
		if answer == "s" {
			break
		}
		if answer != "h" && answer != "s" {
			fmt.Printf("%v is an invalid choice. enter h or s\n", answer)
		}
	}
	return hand
}

// initialDeal sets up the decks of the dealer
// and players at the game's initial state.
func (g *Game) initialDeal() {
	for i := 0; i < 2; i++ {
		for j := 0; j < g.numPlayers; j++ {
			// Name and initialize all players the first time around.
			// Deal them a card both times.
			if i == 0 {
				g.players = append(g.players, player{name: fmt.Sprintf("p_%d", j+1)})
			}
			card := oneOffTheTop(&g.gameDeck)
			g.players[j].hand = append(
				g.players[j].hand,
				card,
			)
			fmt.Printf("dealing %v a %v...\n", g.players[j].name, card)
			time.Sleep(time.Millisecond * 300)
		}
		if i == 0 {
			// Deal the dealer's visible card.
			card := oneOffTheTop(&g.gameDeck)
			g.dealer.hand = append(g.dealer.hand, card)

			fmt.Printf("dealing dealer a %v...\n", card)
			time.Sleep(time.Millisecond * 300)
		}
	}
	// Deal the dealer's hidden card.
	g.dealer.hidden = oneOffTheTop(&g.gameDeck)
	fmt.Println("dealing dealer's hidden card...")
}

// (taking aces into account) ---> if there is an ace, check to see if would be closer to 21 if an ace was instead 11 (even if you have 4, will never have to consider the possibility of more than 1 making it closer to 21)
func handTotal(hand []Card) (total int) {
	var foundAce bool
	for _, c := range hand {
		num := c.Rank
		if num == 1 {
			foundAce = true
		}
		if num > 10 {
			num = 10
		}
		total += int(num)
	}
	// if an Ace is found, check if total would be closer
	// to 21 with 10 added
	if foundAce {
		possible := total + 10
		if 21-possible > -1 && (21-possible < 21-total) {
			total = possible
		}
	}
	return total
}

// hasBlackjack is a helper to check if someone has blackjack.
func hasBlackjack(hand []Card) bool {
	total := handTotal(hand)
	if total == 21 && len(hand) == 2 {
		return true
	}
	return false
}

// oneOffTheTop is a helper to remove a card from the
func oneOffTheTop(d *Deck) (c Card) {
	c, *d = (*d)[0], (*d)[1:]
	return c
}

// printPlayerHands is a helper to print all player hands.
func (g *Game) printPlayerHands() {
	fmt.Println("------------------------------------------")
	for i, p := range g.players {
		fmt.Printf("* * * Player %d * * *\n", i+1)
		for _, c := range p.hand {
			fmt.Println(c)
		}
		fmt.Println("")
	}

	fmt.Println("* * * Dealer * * *")
	for _, c := range g.dealer.hand {
		fmt.Println(c)
	}
	fmt.Println("------------------------------------------")
}

// printFinalScores is a helper to print final scores.
func (g *Game) printFinalScores() {
	fmt.Println("\n=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+")
	fmt.Println("FINAL SCORES")
	fmt.Println("=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+")

	g.printPlayerHands()

	fmt.Println("winner(s):")

	dealerScore := handTotal(g.dealer.hand)

	// If dealer is below 21, everyone above the dealer who
	// hasn't busted wins. If no one who hasn't reached 21 yet
	// is above the dealer, the dealer wins.
	if dealerScore < 21 {
		g.printPlayersInRange(20, 22)
		if 0 == g.printPlayersInRange(dealerScore, 21) {
			fmt.Println("dealer")
		}
	}

	// If dealer has 21, either they got blackjack or didn't.
	// If they did, players with 21 can tie. If not, only past winners.
	if dealerScore == 21 {
		fmt.Println("dealer")
		if hasBlackjack(g.dealer.hand) {
			g.printPlayersInRange(0, 22)
		}
		g.printPlayersInRange(20, 22)
	}

	// If dealer busted, everyone who hasn't busted wins.
	if dealerScore > 21 {
		g.printPlayersInRange(0, 22)
	}
}

// printPlayersInRange prints all players between the lower and
// upper bounds
func (g *Game) printPlayersInRange(lowBound, upBound int) (inRange int) {
	for _, p := range g.players {
		score := handTotal(p.hand)
		if score > lowBound && score < upBound {
			fmt.Println(p.name)
			inRange++
		}
	}
	return inRange
}
