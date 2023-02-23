package akizuki

import (
	"net/http"
	"net/url"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const defaultCatalogURL = "https://akizukidenshi.com/catalog/e/enewall_dL/"

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
	findItemTable(node, &urls)
	return urls, nil
}

func findItemTable(node *html.Node, dst *[]string) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.DataAtom == atom.Table {
				if hasAttr(c, "width", "600") {
					if s := parseItemTable(c); s != "" {
						if url, err := url.JoinPath(defaultBasePath, s); err == nil {
							*dst = append(*dst, url)
						}
					}
				}
			}
			findItemTable(c, dst)
		}
	}
}

func parseItemTable(node *html.Node) string {
	td := getFirstByClassName(node, "cart_tdl")
	if td == nil {
		return ""
	}
	href := findLink(td)
	return href
}
