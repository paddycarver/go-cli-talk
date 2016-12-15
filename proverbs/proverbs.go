package proverbs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Quote struct {
	ID    string
	Value string `json:"Quote"`
}

func GetProverb(baseURL, id string, h http.Header) (Quote, error) {
	u := strings.TrimRight(baseURL, "/") + "/"
	if id != "" {
		q := url.Values{
			"quote": []string{id},
		}
		u = u + "?" + q.Encode()
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return Quote{}, err
	}
	req.Header.Set("Accept", "application/json")
	for k, v := range h {
		req.Header[k] = append(req.Header[k], v...)
	}
	client := &http.Client{Timeout: time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Quote{}, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusBadRequest:
		return Quote{}, errors.New("bad request")
	case http.StatusNotFound:
		return Quote{}, errors.New("that proverb doesn't exist")
	case http.StatusInternalServerError:
		return Quote{}, errors.New("proverbs server ran into a problem, try again later")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Quote{}, err
	}
	var quote Quote
	err = json.Unmarshal(data, &quote)
	return quote, err
}
