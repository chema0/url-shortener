package repository

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/chema0/url-shortener/pkg/utils"
)

// GenerateHash generates a random hash value representing the short form of a given url.
// The returned string hash is the result of concatenating the first 5 runes of the MD5
// hash of the url and 3 random elements from the alphabet.
//
// This approach may not the best, in order to to avoid collisions in a real service with
// a high volume of urls go for a more scalable solution.
// https://blog.codinghorror.com/url-shortening-hashes-in-practice/
func GenerateHash(url string) string {
	h := md5.New()
	hash := h.Sum([]byte(url))
	return utils.Slice(hex.EncodeToString(hash), 5) + generateRandStringFromAlphabet(3)
}

var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func generateRandStringFromAlphabet(length int) string {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	hash := make([]rune, length)
	for i := range hash {
		hash[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(hash)
}
