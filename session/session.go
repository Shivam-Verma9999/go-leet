package session

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/Shivam-Verma9999/go-leetcode/config"
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

	var cookies []*http.Cookie

	for _, pair := range strings.Split(cfg.Cookie, ";") {
		pair := strings.Trim(pair, " ")
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)

		if len(parts) == 2 {
			cookie:= http.Cookie{
				Name: parts[0],
				Value: parts[1],
			}
			cookies = append(cookies, &cookie)
		}

	}

	jar.SetCookies(u, cookies)

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
