package lacework

import (
	"math/rand"
	"time"
)

var (
	letters  = []rune("=,.@:/-abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	rand_gen = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func randString(lenght int) string {
	r := make([]rune, lenght)
	for i := range r {
		r[i] = letters[rand_gen.Intn(len(letters))]
	}
	return string(r)
}
