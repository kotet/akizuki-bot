package akizuki

import "io"

type Item struct {
	Url      string
	Name     string
	Price    string
	ItemCode string
	Images   []io.ReadCloser
}
