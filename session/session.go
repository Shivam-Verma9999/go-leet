package session

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/Shivam-Verma9999/go-leetcode/config"
	"github.com/Shivam-Verma9999/go-leetcode/constants"
)

type Session struct {
	Client    *http.Client
	CSRFToken string
	BaseUrl   string
}

func New(cfg *config.Config) (*Session, error) {
	jar, err := cookiejar.New(nil)

	// pre-load cookies
	u, _ := url.Parse(cfg.BaseUrl)
	jar.SetCookies(u, []*http.Cookie{
		{
			Name:  "cookie",
			Value: cfg.Cookie,
		},
		{
			Name:  constants.CSRFHEADER,
			Value: cfg.CSRFToken,
		},
	})

	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	session := &Session{
		Client:    client,
		CSRFToken: cfg.CSRFToken,
		BaseUrl:   cfg.BaseUrl,
	}
	return session, nil
}
