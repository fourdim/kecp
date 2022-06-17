package kecpfakews

import (
	"math/rand"
	"time"
)

func MathRandGen() int {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	return randGen.Intn(20)
}

func MathRandLongTimeGen() time.Duration {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Duration(randGen.Intn(2000)+1000) * time.Millisecond
}

func MathRandShortTimeGen() time.Duration {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Duration(randGen.Intn(200)+100) * time.Millisecond
}
