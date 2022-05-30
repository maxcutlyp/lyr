package query

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
    "bytes"

	"github.com/PuerkitoBio/goquery"
    "golang.org/x/net/html"
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

        // inspired by goquery's Selection.Text() method but extended to support bold, italics, and newlines
        // https://github.com/PuerkitoBio/goquery/blob/master/property.go#L62
        var buf bytes.Buffer
        var parse_node func(*html.Node)
        parse_node = func(node *html.Node) {
            switch node.Type {
                case html.TextNode:
                    buf.WriteString(node.Data)
                    return
                case html.ElementNode:
                    var open_esc_code string
                    var close_esc_code string

                    switch node.Data {
                        case "h2": // there's a h2 at the top with "$(songtitle) lyrics" that we want to exclude
                            return
                        case "br":
                            buf.WriteString("\n")
                            return
                        case "div":
                            for _, attr := range node.Attr {
                                switch attr.Key {
                                    case "class":
                                        if strings.HasPrefix(attr.Val, "RightSidebar__Container") ||
                                           strings.HasPrefix(attr.Val, "Lyrics__Footer") {
                                            return
                                        }
                                    case "data-exclude-from-selection":
                                        buf.WriteString("\n")
                                }
                            }
                        case "b":
                            open_esc_code = "\033[1m"
                            close_esc_code = "\033[22m"
                        case "i":
                            open_esc_code = "\033[3m"
                            close_esc_code = "\033[23m"
                    }

                    if node.FirstChild != nil {
                        buf.WriteString(open_esc_code)
                        for child := node.FirstChild; child != nil; child = child.NextSibling {
                            parse_node(child)
                        }
                        buf.WriteString(close_esc_code)
                    }
            }
        }

        sel := doc.Find("#lyrics-root")
        for _, node := range sel.Nodes {
            parse_node(node)
        }

		return strings.Trim(buf.String(), "\n "), nil
	}

	for lyrics == "" {
		lyrics, err = attempt()
		if err != nil {
			return "", err
		}
	}
	return lyrics, nil
}
