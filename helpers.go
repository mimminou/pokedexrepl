package main

import (
	"fmt"
	"math/rand"
	"time"
)

func calculateCatchProbability(pokemonXP int32) bool {
	//this is assuming lowest pokemon XP is 1, highest is 255, used 90 to get increments value per 1xp from 10% to 100%. because I decided 10% is the base catch XP, you will have at least 10% chance to catch any given pokemon regardless of other factors

	var incrementChance float32 = 90.0 / 608.0
	catchProb := 100 - (incrementChance * float32(pokemonXP))
	randomness := rand.New(rand.NewSource(time.Now().UnixNano()))
	missChance := (0.1 + randomness.Float32()*0.9) * 100.0
	fmt.Println(fmt.Sprintf("Pokemon Base XP : %d", pokemonXP))
	fmt.Println(fmt.Sprintf("catch probability : %d%%", int(catchProb)))

	if missChance < catchProb {
		return true
	}
	return false
}
