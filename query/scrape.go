package query

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (data queryData) ScrapeLyrics() (lyrics string, err error) {

	attempt := func() (string, error) {

		res, err := http.Get("https://genius.com" + data.Path)
		if err != nil {
			return "", err
		}
		if res.StatusCode != http.StatusOK {
			return "", errors.New(fmt.Sprintf(
				"Received unexpected status code %d\n", res.StatusCode,
			))
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return "", err
		}

		return strings.Trim(
			doc.Find(".lyrics").
				First().
				Text(),
			"\n "), nil
	}

	for lyrics == "" {
		lyrics, err = attempt()
		if err != nil {
			return "", err
		}
	}
	return lyrics, nil
}
