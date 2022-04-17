package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	sim := NewCrapsSim(200, 100)
	sim.SetBettingStrategy(BettingStrategyPassLine, 15)
	sim.SetBettingStrategy(BettingStrategyOdds, 15)
	sim.SetBettingStrategy(BettingStrategyPlaceBet6, 6)
	sim.SetBettingStrategy(BettingStrategyPlaceBet8, 6)
	sim.Start()
}

func generateRandomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}

type Craps struct {
	isPuckOn          bool
	Bankroll          int
	iteration         int
	MaxIteration      int
	round             int
	Die1              int
	Die2              int
	PointValue        int
	BettingStrategies *[]BettingStrategy
	board             *Board
}

const (
	BettingStrategyPlaceBet4    = 4
	BettingStrategyPlaceBet5    = 5
	BettingStrategyPlaceBet6    = 6
	BettingStrategyPlaceBet8    = 8
	BettingStrategyPlaceBet9    = 9
	BettingStrategyPlaceBet10   = 10
	BettingStrategyFieldBet     = 20
	BettingStrategyPassLine     = 21
	BettingStrategyDontPassLine = 22
	BettingStrategyOdds         = 23
)

type BettingStrategy struct {
	StrategyType        int
	Bet                 int
	TakeDownBetAfterWin bool
}

type Board struct {
	PlaceBet4      int
	PlaceBet5      int
	PlaceBet6      int
	PlaceBet8      int
	PlaceBet9      int
	PlaceBet10     int
	PassLineBet    int
	PassLineOdds   int
	DontPassBarBet int
	FieldBet       int
	ComeBet        int
}

func NewCrapsSim(startBalance int, maxIterations int) *Craps {
	return &Craps{
		Bankroll:     startBalance,
		MaxIteration: maxIterations,
		board:        &Board{},
	}
}

func (c *Craps) rollDie() int {
	return generateRandomNumber(1, 6)
}

func (c *Craps) rollDice() {
	c.Die1 = c.rollDie()
	c.Die2 = c.rollDie()
}

func (c *Craps) getDiceValue() int {
	return c.Die1 + c.Die2
}

func (c *Craps) SetBettingStrategy(StrategyType int, bet int) {
	strat := BettingStrategy{
		StrategyType: StrategyType,
		Bet:          bet,
	}

	if c.BettingStrategies == nil {
		c.BettingStrategies = &[]BettingStrategy{}
	}
	*c.BettingStrategies = append(*c.BettingStrategies, strat)
}

func (c *Craps) Start() {
	c.round = 1
	c.iteration = 1

	fmt.Printf("Initial bankroll is set to $%d\n", c.Bankroll)
	c.setPassLineBet()

	for c.Bankroll > 0 && c.iteration <= c.MaxIteration {
		c.rollDice()
		c.printRoundIteration()
		fmt.Printf("Dice rolled %d\n", c.getDiceValue())

		if !c.isPuckOn {
			if c.getDiceValue() == 2 || c.getDiceValue() == 3 || c.getDiceValue() == 12 {
				// craps
				c.printRoundIteration()
				fmt.Printf("Rolled craps. End of round.\n")
				c.board.PassLineBet = 0
				c.round += 1
				c.setPassLineBet()
				c.printBankRoll()
			} else if c.getDiceValue() == 7 || c.getDiceValue() == 11 {
				// win
				c.printRoundIteration()
				fmt.Printf("Win Come Out roll. End of round.\n")
				c.Bankroll += c.board.PassLineBet // keep original bet on field and return winnings to bank roll
				c.printBankRoll()
				c.round += 1
				c.setPassLineBet()
				c.printBankRoll()
			} else {
				// point is set
				c.PointValue = c.getDiceValue()
				c.isPuckOn = true
				c.setBoard()
				c.printRoundIteration()
				fmt.Printf("Point is set at %d\n", c.getDiceValue())
			}
		} else if c.isPuckOn {
			if c.getDiceValue() == 7 {
				// seven-out (lose)
				c.printRoundIteration()
				fmt.Printf("Rolled seven-out (loss). End of round.\n")
				c.calculateSevenOutLoss()
				c.isPuckOn = false
				c.round += 1
			}

			var currentBetStrategy *BettingStrategy = nil

			currentBetStrategy = c.isBettingStrategyExists(BettingStrategyPlaceBet4)
			if c.getDiceValue() == 4 && currentBetStrategy != nil {
				payoutMultiplier := currentBetStrategy.Bet / 5
				payout := (payoutMultiplier * 4) + currentBetStrategy.Bet
				c.Bankroll += payout
				c.printRoundIteration()
				fmt.Printf("Won on point 4. Payout is $%d. Bankroll increased to $%d.\n", payout, c.Bankroll)

				if currentBetStrategy.TakeDownBetAfterWin {
					c.Bankroll += currentBetStrategy.Bet
					c.board.PlaceBet4 = 0
					c.printRoundIteration()
					fmt.Printf("Taking down bet on point 4. Bankroll increased to $%d.\n", c.Bankroll)
				}
			}

			currentBetStrategy = c.isBettingStrategyExists(BettingStrategyPlaceBet5)
			if c.getDiceValue() == 5 && currentBetStrategy != nil {
				payoutMultiplier := currentBetStrategy.Bet / 5
				payout := (payoutMultiplier * 2) + currentBetStrategy.Bet
				c.Bankroll += payout
				c.printRoundIteration()
				fmt.Printf("Won on point 5. Payout is $%d. Bankroll increased to $%d.\n", payout, c.Bankroll)

				if currentBetStrategy.TakeDownBetAfterWin {
					c.Bankroll += currentBetStrategy.Bet
					c.board.PlaceBet5 = 0
					c.printRoundIteration()
					fmt.Printf("Taking down bet on point 5. Bankroll increased to $%d.\n", c.Bankroll)
				}
			}

			currentBetStrategy = c.isBettingStrategyExists(BettingStrategyPlaceBet6)
			if c.getDiceValue() == 6 && currentBetStrategy != nil {
				payoutMultiplier := currentBetStrategy.Bet / 6
				payout := (payoutMultiplier * 1) + currentBetStrategy.Bet
				c.Bankroll += payout
				c.printRoundIteration()
				fmt.Printf("Won on point 6. Payout is $%d. Bankroll increased to $%d.\n", payout, c.Bankroll)

				if currentBetStrategy.TakeDownBetAfterWin {
					c.Bankroll += currentBetStrategy.Bet
					c.board.PlaceBet6 = 0
					c.printRoundIteration()
					fmt.Printf("Taking down bet on point 6. Bankroll increased to $%d.\n", c.Bankroll)
				}
			}

			currentBetStrategy = c.isBettingStrategyExists(BettingStrategyPlaceBet8)
			if c.getDiceValue() == 8 && currentBetStrategy != nil {
				payoutMultiplier := currentBetStrategy.Bet / 6
				payout := (payoutMultiplier * 1) + currentBetStrategy.Bet
				c.Bankroll += payout
				c.printRoundIteration()
				fmt.Printf("Won on point 8. Payout is $%d. Bankroll increased to $%d.\n", payout, c.Bankroll)

				if currentBetStrategy.TakeDownBetAfterWin {
					c.Bankroll += currentBetStrategy.Bet
					c.board.PlaceBet8 = 0
					c.printRoundIteration()
					fmt.Printf("Taking down bet on point 8. Bankroll increased to $%d.\n", c.Bankroll)
				}
			}

			currentBetStrategy = c.isBettingStrategyExists(BettingStrategyPlaceBet9)
			if c.getDiceValue() == 9 && currentBetStrategy != nil {
				payoutMultiplier := currentBetStrategy.Bet / 5
				payout := (payoutMultiplier * 2) + currentBetStrategy.Bet
				c.Bankroll += payout
				c.printRoundIteration()
				fmt.Printf("Won on point 9. Payout is $%d. Bankroll increased to $%d.\n", payout, c.Bankroll)

				if currentBetStrategy.TakeDownBetAfterWin {
					c.Bankroll += currentBetStrategy.Bet
					c.board.PlaceBet9 = 0
					c.printRoundIteration()
					fmt.Printf("Taking down bet on point 9. Bankroll increased to $%d.\n", c.Bankroll)
				}
			}

			currentBetStrategy = c.isBettingStrategyExists(BettingStrategyPlaceBet10)
			if c.getDiceValue() == 10 && currentBetStrategy != nil {
				payoutMultiplier := currentBetStrategy.Bet / 5
				payout := (payoutMultiplier * 4) + currentBetStrategy.Bet
				c.Bankroll += payout
				c.printRoundIteration()
				fmt.Printf("Won on point 10. Payout is $%d. Bankroll increased to $%d.\n", payout, c.Bankroll)

				if currentBetStrategy.TakeDownBetAfterWin {
					c.Bankroll += currentBetStrategy.Bet
					c.board.PlaceBet10 = 0
					c.printRoundIteration()
					fmt.Printf("Taking down bet on point 10. Bankroll increased to $%d.\n", c.Bankroll)
				}
			}

			if c.getDiceValue() == c.PointValue {
				// win
				c.printRoundIteration()
				fmt.Printf("Rolled point value (win!). End of round.\n")
				c.calculatePassLineAndOddsPayout()
				c.isPuckOn = false
				c.PointValue = 0
				c.round += 1
			}
		}

		c.iteration += 1
	}

	fmt.Printf("Simulation is complete. Bankroll is $%d\n", c.Bankroll)
}

func (c *Craps) isBettingStrategyExists(bettingStrategy int) *BettingStrategy {
	if *c.BettingStrategies != nil {
		for _, strategy := range *c.BettingStrategies {
			if strategy.StrategyType == bettingStrategy {
				return &strategy
			}
		}
	}

	return nil
}

func (c *Craps) setPassLineBet() {
	passLine := c.isBettingStrategyExists(BettingStrategyPassLine)
	c.Bankroll -= passLine.Bet
	c.board.PassLineBet = passLine.Bet

	c.printRoundIteration()
	fmt.Printf("Added $%d to pass line bet.\n", passLine.Bet)
}

func (c *Craps) setBoard() {
	for _, strategy := range *c.BettingStrategies {
		c.Bankroll -= strategy.Bet

		c.printRoundIteration()
		switch strategy.StrategyType {
		case BettingStrategyPlaceBet4:
			fmt.Printf("Added $%d to place bet 4.\n", strategy.Bet)
			c.board.PlaceBet4 = strategy.Bet
		case BettingStrategyPlaceBet5:
			fmt.Printf("Added $%d to place bet 5.\n", strategy.Bet)
			c.board.PlaceBet5 = strategy.Bet
		case BettingStrategyPlaceBet6:
			fmt.Printf("Added $%d to place bet 6.\n", strategy.Bet)
			c.board.PlaceBet6 = strategy.Bet
		case BettingStrategyPlaceBet8:
			fmt.Printf("Added $%d to place bet 8.\n", strategy.Bet)
			c.board.PlaceBet8 = strategy.Bet
		case BettingStrategyPlaceBet9:
			fmt.Printf("Added $%d to place bet 9.\n", strategy.Bet)
			c.board.PlaceBet9 = strategy.Bet
		case BettingStrategyPlaceBet10:
			fmt.Printf("Added $%d to place bet 10.\n", strategy.Bet)
			c.board.PlaceBet10 = strategy.Bet
		case BettingStrategyOdds:
			fmt.Printf("Added $%d to pass line odds bet.\n", strategy.Bet)
			c.board.PassLineOdds = strategy.Bet
		}

		c.printBankRoll()
	}
}

func (c *Craps) calculatePassLineAndOddsPayout() {
	passLine := c.isBettingStrategyExists(BettingStrategyPassLine)
	if passLine != nil {
		c.Bankroll += passLine.Bet * 2
		c.printRoundIteration()
		fmt.Printf("Pay out pass line. Bankroll has increased to $%d\n", c.Bankroll)
	}

	passLineOdds := c.isBettingStrategyExists(BettingStrategyOdds)
	if passLineOdds != nil {
		if c.getDiceValue() == 4 || c.getDiceValue() == 10 {
			// 2 to 1 odds
			c.Bankroll += (passLineOdds.Bet * 2) + passLineOdds.Bet
		} else if c.getDiceValue() == 5 || c.getDiceValue() == 9 {
			// 3 to 2 odds
			c.Bankroll += int(float64(passLineOdds.Bet)*1.5) + passLineOdds.Bet
		} else if c.getDiceValue() == 6 || c.getDiceValue() == 8 {
			// 6 to 5 odds
			c.Bankroll += int(float64(passLineOdds.Bet)*1.2) + passLineOdds.Bet
		}

		c.printRoundIteration()
		fmt.Printf("Pay out odds. Bankroll has increased to $%d\n", c.Bankroll)
	}
}

func (c *Craps) calculateSevenOutLoss() {
	c.board = &Board{}
	c.printBankRoll()
}

func (c *Craps) printRoundIteration() {
	fmt.Printf("[Round %d][Iteration %d] ", c.round, c.iteration)
}

func (c *Craps) printBankRoll() {
	c.printRoundIteration()
	fmt.Printf("Bankroll is set to $%d\n", c.Bankroll)
}
