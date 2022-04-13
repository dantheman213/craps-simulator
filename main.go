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
	Bankroll          float64
	iteration         int
	MaxIteration      int
	round             int
	Die1              int
	Die2              int
	PointValue        int
	BettingStrategies *[]BettingStrategy
	placeBetHistory   *[]int
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
	Bet                 float64
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
	DontPassBarBet int
	FieldBet       int
	ComeBet        int
}

func NewCrapsSim(startBalance float64, maxIterations int) *Craps {
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

func (c *Craps) SetBettingStrategy(StrategyType int, bet float64) {
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

	fmt.Printf("Initial bankroll is set to $%f\n", c.Bankroll)
	for c.Bankroll > 0 && c.iteration <= c.MaxIteration {
		c.rollDice()
		fmt.Printf("[Round %d][Iteration %d] Dice rolled %d\n", c.round, c.iteration, c.getDiceValue())

		if !c.isPuckOn {
			if c.getDiceValue() == 2 || c.getDiceValue() == 3 || c.getDiceValue() == 12 {
				// craps
				fmt.Printf("[Round %d][Iteration %d] Rolled craps. End of round.\n", c.round, c.iteration)
				c.round += 1
			} else if c.getDiceValue() == 7 || c.getDiceValue() == 11 {
				// win
				fmt.Printf("[Round %d][Iteration %d] Win Come Out roll. End of round.\n", c.round, c.iteration)
				c.round += 1
			} else {
				// point is set
				c.PointValue = c.getDiceValue()
				c.isPuckOn = true
				fmt.Printf("[Round %d][Iteration %d] Point is set at %d\n", c.round, c.iteration, c.getDiceValue())
			}
		} else if c.isPuckOn {
			if c.getDiceValue() == 7 {
				// seven-out (lose)
				fmt.Printf("[Round %d][Iteration %d] Rolled seven-out (loss). End of round.\n", c.round, c.iteration)

				c.calculateSevenOutLoss()

				c.isPuckOn = false
				c.round += 1
				c.placeBetHistory = &[]int{}
			}

			var currentBetStrategy *BettingStrategy = nil

			currentBetStrategy = c.isBettingStrategyExists(BettingStrategyPlaceBet4)
			if c.getDiceValue() == 4 && currentBetStrategy != nil {
				payoutMultiplier := currentBetStrategy.Bet / 5
				payout := (payoutMultiplier * 4) + currentBetStrategy.Bet
				c.Bankroll += payout
				fmt.Printf("[Round %d][Iteration %d] Won on point 4. Payout is %f. Bankroll increased to $%f.\n", c.round, c.iteration, payout, c.Bankroll)

				if currentBetStrategy.TakeDownBetAfterWin {
					c.Bankroll += currentBetStrategy.Bet
					fmt.Printf("[Round %d][Iteration %d] Taking down bet on point 4. Bankroll increased to $%f.\n", c.round, c.iteration, c.Bankroll)
				}
			}

			if c.getDiceValue() == c.PointValue {
				// win

				fmt.Printf("[Round %d][Iteration %d] Rolled point value (win!). End of round.\n", c.round, c.iteration)
				c.isPuckOn = false
				c.PointValue = 0
				c.round += 1
				c.placeBetHistory = &[]int{}
			}
		}

		c.iteration += 1
	}
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

func (c *Craps) addPlaceBetHistory(value int) {
	if c.placeBetHistory == nil {
		c.placeBetHistory = &[]int{}
	}

	*c.placeBetHistory = append(*c.placeBetHistory, value)
}

func (c *Craps) calculateSevenOutLoss() {

}
