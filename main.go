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
	sim.SetBettingStrategy(BettingStrategyPoint6, 6)
	sim.SetBettingStrategy(BettingStrategyPoint8, 6)
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
}

const (
	BettingStrategyPoint4       = 4
	BettingStrategyPoint5       = 5
	BettingStrategyPoint6       = 6
	BettingStrategyPoint8       = 8
	BettingStrategyPoint9       = 9
	BettingStrategyPoint10      = 10
	BettingStrategyFieldBet     = 20
	BettingStrategyPassLine     = 21
	BettingStrategyDontPassLine = 22
	BettingStrategyOdds         = 23
)

type BettingStrategy struct {
	StrategyType int
	Bet          float64
}

func NewCrapsSim(startBalance float64, maxIterations int) *Craps {
	return &Craps{
		Bankroll:     startBalance,
		MaxIteration: maxIterations,
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
				c.isPuckOn = false
				c.round += 1
			} else if c.getDiceValue() == c.PointValue {
				// win

				fmt.Printf("[Round %d][Iteration %d] Rolled point value (win!). End of round.\n", c.round, c.iteration)
				c.isPuckOn = false
				c.PointValue = 0
				c.round += 1
			} else {
				// TBA

			}
		}

		c.iteration += 1
	}
}
