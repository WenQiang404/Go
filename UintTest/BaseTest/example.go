package main

import (
	"math/rand"
)

var ServerIndex [10]int

func InitServerIndex() {
	for i := 0; i < 10; i++ {
		ServerIndex[i] = 1 + 100
	}
}

func Select() int {
	return ServerIndex[rand.Intn(10)]
}
