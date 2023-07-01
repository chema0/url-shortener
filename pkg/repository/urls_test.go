package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/chema0/url-shortener/config"
	"github.com/chema0/url-shortener/pkg/db"
	"github.com/chema0/url-shortener/pkg/repository"
)

func TestCreateURL(t *testing.T) {
	// ctx := db.TestContext(t)
	db.TestContext(t)
	cfg := config.LoadConfig()

	url := repository.URL{
		URL:       "https://www.twitch.tv/directory/game/Software%20and%20Game%20Development",
		ExpiredAt: time.Now(),
		OwnerID:   "test",
	}

	urlsRepository := repository.NewURLs(db.Global, &cfg.Database)
	newUrl, err := urlsRepository.CreateURL(url)
	if err != nil {
		t.Fatal("failed to load config")
	}

	fmt.Println(newUrl)
}
