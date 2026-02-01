package utils

import (
	"math/rand"
	"time"
)

var randomNumber = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateRandomNumber(minVal, maxVal int) int {
	return randomNumber.Intn(maxVal-minVal) + minVal
}
