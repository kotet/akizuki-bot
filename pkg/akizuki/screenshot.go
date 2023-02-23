package akizuki

import (
	"bytes"
	"context"
	"io"
	"log"

	screenshots "github.com/screenshotone/gosdk"
)

type ScreenShotOneConfig struct {
	AccessKey string `yaml:"access"`
	SecretKey string `yaml:"secret"`
}

func ScreenShotOneTakeScreenShot(config ScreenShotOneConfig) func(url string) (io.Reader, error) {
	return func(pageURL string) (io.Reader, error) {
		client, err := screenshots.NewClient(config.AccessKey, config.SecretKey)
		if err != nil {
			return nil, err
		}
		options := screenshots.NewTakeOptions(pageURL).Format("png")
		image, _, err := client.Take(context.TODO(), options)
		if err != nil {
			return nil, err
		}
		log.Println("Screenshot taken")
		return io.NopCloser(bytes.NewReader(image)), nil
	}
}
