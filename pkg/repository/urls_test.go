package repository_test

import (
	"testing"

	"github.com/chema0/url-shortener/config"
	"github.com/chema0/url-shortener/pkg/db"
	"github.com/chema0/url-shortener/pkg/repository"
	"github.com/stretchr/testify/assert"
)

func setupTest(tb testing.TB) repository.URLs {
	db.TestContext(tb)
	cfg := config.LoadConfig()
	return repository.NewURLs(db.Global, &cfg.Database)
}

func TestCreateURL(t *testing.T) {
	repo := setupTest(t)

	u := newURL()

	url, err := repo.CreateURL(u)
	assert.NoError(t, err)

	assert.Equal(t, url.URL, u.URL)
	assert.Equal(t, url.OwnerID, u.OwnerID)
	assert.NotEmpty(t, url.Hash)
	assert.NotEmpty(t, url.CreatedAt)
}

func TestListURLs(t *testing.T) {
	repo := setupTest(t)

	u1, _ := repo.CreateURL(newURL())
	u2, _ := repo.CreateURL(newURL())

	urls, err := repo.ListURLs("test")
	assert.NoError(t, err)
	assert.Len(t, urls, 2)
	assert.Equal(t, u1, urls[0])
	assert.Equal(t, u2, urls[1])
}

func TestGetURLByHash(t *testing.T) {
	repo := setupTest(t)

	u, _ := repo.CreateURL(newURL())

	url, err := repo.GetURLByHash(u.Hash)
	assert.NoError(t, err)
	assert.Equal(t, u, url)
}

func TestDeleteURL(t *testing.T) {
	repo := setupTest(t)
	u, _ := repo.CreateURL(newURL())

	err := repo.DeleteURL(u.ID)
	assert.NoError(t, err)

	_, err = repo.GetURLByHash(u.Hash)
	assert.Error(t, err)
}

func newURL() repository.URL {
	return repository.URL{
		URL:     "https://www.twitch.tv/directory/game/Software%20and%20Game%20Development",
		OwnerID: "test",
	}
}
