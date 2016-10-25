package proverbs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Quote{}, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Quote{}, err
	}
	var quote Quote
	err = json.Unmarshal(data, &quote)
	return quote, err
}
