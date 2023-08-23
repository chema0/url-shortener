package repository

import (
	"crypto/md5"
	"io"
)

func generateHash(url string) string {
	h := md5.New()
	io.WriteString(h, url)
	return string(h.Sum(nil)) + 
}

var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
// func generateRandHash(length int) string {
// 	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	hash := make([]rune, length)
// 	for i := range hash {
// 		hash[i] = alphabet[rand.Intn(len(alphabet))]
// 	}
// 	return string(hash)
// }
