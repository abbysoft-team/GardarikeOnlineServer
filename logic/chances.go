package logic

import (
	"math/rand"
	"time"
)

const (
	PopulationGrownEventChance = 2.0
	TreeGrownEventChance       = 2.0
)

// checkRandomEventHappened - check if the random event of 'chance' percent freq has happened
// returns true if the event has happened and false otherwise
func CheckRandomEventHappened(chance int) bool {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int63n(10e10)
	chance *= 10e8

	if n >= int64(chance) {
		return false
	}

	return true
}
