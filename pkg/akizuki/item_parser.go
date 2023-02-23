package akizuki

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
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

		shiftJISReader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())

		node, err := html.Parse(shiftJISReader)
		if err != nil {
			return nil, err
		}
		// retrieve image
		imgLink := getFirstById(node, "imglink")
		if imgLink == nil {
			return nil, errors.New("failed to parse")
		}
		imgURL := getAttr(imgLink, "href")

		fullImgURL, err := url.JoinPath(defaultBasePath, imgURL)
		if err != nil {
			return nil, err
		}
		resp, err = http.Get(fullImgURL)
		if err != nil {
			resp.Body.Close()
			return nil, err
		}

		// parse order info
		order := getFirstByClassName(node, "order_g")
		if order == nil {
			return nil, errors.New("failed to parse")
		}

		name, price, itemCode, err := parseOrder(order)
		if err != nil {
			return nil, err
		}

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

func parseOrder(order *html.Node) (name string, price string, itemCode string, err error) {
	title := strings.Split(getFirstText(order), "\u3000")
	itemCode = strings.Trim(title[0], "[]")
	name = strings.Trim(strings.Join(title[1:], " "), " ")

	table := getFirstByAtom(order, atom.Table)
	if table == nil {
		err = errors.New("failed to parse")
		return
	}
	tr := getFirstByAtom(table, atom.Tr)
	if tr == nil {
		err = errors.New("failed to parse")
		return
	}
	td := getFirstByAtom(table, atom.Td)
	if td == nil {
		err = errors.New("failed to parse")
		return
	}
	price = strings.ReplaceAll(getText(td), "\t", "")
	price = strings.ReplaceAll(price, "\n", " ")
	price = strings.Trim(price, " ")
	return
}
