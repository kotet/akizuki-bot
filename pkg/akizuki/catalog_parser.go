package akizuki

import (
	"net/http"
	"net/url"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const defaultCatalogURL = "https://akizukidenshi.com/catalog/e/enewall_dL/"

// retrieve catalog page and find new item urls
func defaultCatalogParser() ([]string, error) {
	resp, err := http.Get(defaultCatalogURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	urls := []string{}
	findItemList(node, &urls)
	return urls, nil
}

func findItemList(node *html.Node, dst *[]string) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.DataAtom == atom.Ul {
				if hasAttr(c, "class", "block-cart-i--items") {
					parseItemList(c, dst)
					return
				}
			}
			findItemList(c, dst)
		}
	}
}

func parseItemList(node *html.Node, dst *[]string) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.DataAtom == atom.Dl {
				url := parseItem(c)
				if url != "" {
					*dst = append(*dst, url)
				}
			}
			parseItemList(c, dst)
		}
	}
}

func parseItem(node *html.Node) string {
	a := getFirstByClassName(node, "js-enhanced-ecommerce-goods-name")
	if a == nil {
		return ""
	}
	href := getAttr(a, "href")
	u, err := url.JoinPath(defaultBasePath, href)
	if err != nil {
		return ""
	}
	return u
}
