package utils

import (
	"math/rand"
	"time"
)

var Source = rand.NewSource(time.Now().UnixNano())
var RandomNumber = rand.New(Source)
