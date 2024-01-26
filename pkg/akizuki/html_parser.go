package akizuki

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func getFirstByClassName(node *html.Node, classname string) *html.Node {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			class := getAttr(c, "class")
			classes := strings.Split(class, " ")
			for _, cls := range classes {
				if cls == classname {
					return c
				}
			}
			if n := getFirstByClassName(c, classname); n != nil {
				return n
			}
		}
	}
	return nil
}

func getFirstByAtom(node *html.Node, atom atom.Atom) *html.Node {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.DataAtom == atom {
				return c
			}
			if n := getFirstByAtom(c, atom); n != nil {
				return n
			}
		}
	}
	return nil
}

func getFirstById(node *html.Node, id string) *html.Node {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if hasAttr(c, "id", id) {
				return c
			}
			if n := getFirstById(c, id); n != nil {
				return n
			}
		}
	}
	return nil
}

func getFirstText(node *html.Node) string {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			return c.Data
		}
	}
	return ""
}

func getText(node *html.Node) string {
	builder := strings.Builder{}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			builder.WriteString(c.Data)
		}
		if c.Type == html.ElementNode {
			builder.WriteString(getText(c))
		}
	}
	return builder.String()
}

func hasAttr(node *html.Node, key string, val string) bool {
	for _, attr := range node.Attr {
		if attr.Key == key && attr.Val == val {
			return true
		}
	}
	return false
}

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func findLink(node *html.Node) string {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.DataAtom == atom.A {
				for _, attr := range c.Attr {
					if attr.Key == "href" {
						return attr.Val
					}
				}
			}
			if s := findLink(node); s != "" {
				return s
			}
		}
	}
	return ""
}
