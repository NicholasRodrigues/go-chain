package main

import (
	"math/rand"
	"strconv"
	"time"
)

//  T is the target value when using a proof of work algorithm
func GenerateT() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	prefix := rand.Intn(90) + 10

	suffix := rand.Intn(10000000)

	T := strconv.Itoa(prefix) + strconv.Itoa(suffix)

	return T
}
