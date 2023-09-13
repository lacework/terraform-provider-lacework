package lacework

import (
	"math/rand"
	"time"
)

var letters = []rune("=,.@:/-abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randString(lenght int) string {
	var rand_gen = rand.New(rand.NewSource(time.Now().UnixNano()))
	r := make([]rune, lenght)
	for i := range r {
		r[i] = letters[rand_gen.Intn(len(letters))]
	}
	return string(r)
}
