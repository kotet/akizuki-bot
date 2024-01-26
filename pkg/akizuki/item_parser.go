package akizuki

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func defaultParseItem(minWait time.Duration) func(itemURL string) (*Item, error) {
	lastVisit := time.Now()
	return func(itemURL string) (*Item, error) {
		if lastVisit.Add(minWait).After(time.Now()) {
			diff := time.Until(lastVisit.Add(minWait))
			time.Sleep(diff)
		}
		resp, err := http.Get(itemURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		node, err := html.Parse(resp.Body)
		if err != nil {
			return nil, err
		}

		imgURL := findImageURL(node)

		fullImgURL, err := url.JoinPath(defaultBasePath, imgURL)
		if err != nil {
			return nil, err
		}
		resp, err = http.Get(fullImgURL)
		if err != nil {
			resp.Body.Close()
			return nil, err
		}

		name := parseTitle(node)
		price := parsePrice(node)
		itemCode := parseItemCode(node)

		return &Item{
			Url:      itemURL,
			Name:     name,
			Price:    price,
			ItemCode: itemCode,
			Images: []io.ReadCloser{
				resp.Body,
			},
		}, nil
	}
}

func findImageURL(node *html.Node) string {
	d := getFirstByClassName(node, "block-src-l")
	if d == nil {
		return ""
	}
	i := getFirstByAtom(d, atom.Img)
	if i == nil {
		return ""
	}
	return getAttr(i, "data-src")
}

func parseTitle(node *html.Node) string {
	h1 := getFirstByAtom(node, atom.H1)
	if h1 == nil {
		return ""
	}
	return getText(h1)
}

func parseItemCode(node *html.Node) string {
	dd := getFirstById(node, "spec_goods")
	if dd == nil {
		return ""
	}
	return getText(dd)
}

func parsePrice(node *html.Node) string {
	div := getFirstByClassName(node, "block-goods-price--price")
	if div == nil {
		return ""
	}
	builder := strings.Builder{}
	for c := div.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			builder.WriteString(c.Data)
		}
	}
	return strings.Trim(builder.String(), " \n\t")
}
