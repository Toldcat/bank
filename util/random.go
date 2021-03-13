package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) //returns numbers between min - max
}

func RandomString(n int) string {
	//initialise a string builder
	var sb strings.Builder
	//get alphabet size
	k := len(alphabet)

	//loop through passed number in the function
	for i := 0; i < n; i++ {
		//take a random letter from the alphabet
		c := alphabet[rand.Intn(k)]
		//write that letter into the string builder
		sb.WriteByte(c)
	}
	//return the completed string
	return sb.String()
}

//Generate random Name
func RandomName() string {
	return RandomString(6)
}

//generate random amount
func RandomAmount() int64 {
	return RandomInt(0, 1000)
}

//Generate Random currency
func RandomCurrency() string {
	currlist := []string{"USD", "GBP", "EUR"}
	n := len(currlist)

	currency := currlist[rand.Intn(n)]
	return currency

}
