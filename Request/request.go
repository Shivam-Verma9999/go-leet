package request

import (
	"io"
	"net/http"
)

func NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	request.Header.Set("Accept-Language", "en-US,en;q=0.5")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Referer", "https://leetcode.com/problems/two-sum/description/")
	request.Header.Set("Content-Type", "application/json")

	return request, nil

}

func MakeRequest(request *http.Request, client *http.Client) (*http.Response, error) {

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
