package repository

import (
	"math/rand"
	"time"
)

var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func generateHash(length int) string {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	hash := make([]rune, length)
	for i := range hash {
		hash[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(hash)
}
