package repository_test

import (
	"crypto/md5"
	"encoding/hex"
	"testing"
	"unicode/utf8"

	"github.com/chema0/url-shortener/pkg/repository"
	"github.com/chema0/url-shortener/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateHash(t *testing.T) {
	const url = "https://hexdocs.pm/elixir/1.15.4/Kernel.html"
	hash := repository.GenerateHash(url)
	md5Hash := md5.New().Sum([]byte(url))

	assert.Equal(t, utf8.RuneCountInString(hash), 8)
	assert.Equal(t, utils.Slice(hash, 5), utils.Slice(hex.EncodeToString(md5Hash), 5))
}
