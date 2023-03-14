package akizuki

import (
	"io"
	"log"
	"time"
)

const defaultBasePath = "https://akizukidenshi.com"

type AkizukiBot struct {
	parseCatalog   func() ([]string, error)
	detector       NewPageDetector
	parseItem      func(url string) (*Item, error)
	takeScreenShot func(url string) (io.Reader, error)
	format         func(item *Item) string
	toot           func(text string, images []io.Reader) error
}

func NewBot() (*AkizukiBot, error) {
	d, err := NewDefaultDetector(defaultDatabasePath)
	if err != nil {
		return nil, err
	}
	return &AkizukiBot{
		parseCatalog: defaultCatalogParser,
		detector:     d,
		parseItem:    defaultParseItem(time.Second),
	}, nil
}

func (b *AkizukiBot) SetParseCatalog(f func() ([]string, error)) *AkizukiBot {
	b.parseCatalog = f
	return b
}
func (b *AkizukiBot) SetNewPageDetector(d NewPageDetector) *AkizukiBot {
	b.detector = d
	return b
}
func (b *AkizukiBot) SetParseItem(f func(url string) (*Item, error)) *AkizukiBot {
	b.parseItem = f
	return b
}
func (b *AkizukiBot) SetTakeScreenshot(f func(url string) (io.Reader, error)) *AkizukiBot {
	b.takeScreenShot = f
	return b
}
func (b *AkizukiBot) SetFormat(f func(item *Item) string) *AkizukiBot {
	b.format = f
	return b
}
func (b *AkizukiBot) SetToot(f func(text string, images []io.Reader) error) *AkizukiBot {
	b.toot = f
	return b
}

func (b *AkizukiBot) RunOnce() error {
	urls, err := b.parseCatalog()
	if err != nil {
		return err
	}

	newURLs, err := b.detector.NewPages(urls)
	if err != nil {
		return err
	}

	for _, url := range newURLs {
		err := func() error {
			log.Printf("New Item Detected: %v", url)
			item, err := b.parseItem(url)
			if err != nil {
				return err
			}
			for _, itemImage := range item.Images {
				defer itemImage.Close()
			}
			images := []io.Reader{}
			if len(item.Images) > 0 {
				images = append(images, item.Images[0])
			}

			if b.takeScreenShot != nil {
				ss, err := b.takeScreenShot(url)
				if err != nil {
					log.Printf("failed to take screenshot: %v", err.Error())
				} else {
					images = append(images, ss)
				}
			}

			s := b.format(item)
			err = b.toot(s, images)
			if err != nil {
				return err
			}
			err = b.detector.AddPages([]string{url})
			if err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}
