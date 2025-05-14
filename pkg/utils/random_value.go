package utils

import (
	"math/rand"
	"time"
)

var randomNumber = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateRandomNumber(min, max int) int {
	return randomNumber.Intn(max-min) + min
}
