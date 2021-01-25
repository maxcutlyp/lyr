package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/satoqz/lyr/config"
)

type Query struct {
	Raw     string
	Encoded string
}

func New(q string) Query {
	return Query{
		Raw:     q,
		Encoded: url.QueryEscape(q),
	}
}

func (q Query) Search() (data queryResponse, err error) {

	token, err := config.ReadToken()
	if err != nil {
		return data, err
	}

	url := "https://api.genius.com/search?q=" + q.Encoded
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return data, err
	}
	if res.StatusCode != http.StatusOK {
		return data, errors.New(fmt.Sprintf(
			"Received unexpected status code %d\nMessage: %s\n", res.StatusCode, res.Status,
		))
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	return data, err
}
