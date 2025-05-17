package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kotet/akizuki-bot/pkg/akizuki"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Mastodon      akizuki.DefaultTootConfig   `yaml:"mastodon"`
	ScreenshotOne akizuki.ScreenShotOneConfig `yaml:"screenshotone"`
}

func main() {

	configPath := "config.yaml"

	configYAML := func() []byte {
		ret, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalln(err.Error())
		}
		return ret
	}()

	var config Config
	err := yaml.Unmarshal(configYAML, &config)
	if err != nil {
		log.Fatalln(err.Error())
	}

	detector, err := akizuki.NewDefaultDetector("akizuki.json")
	if err != nil {
		log.Fatalln(err.Error())
	}

	bot, err := akizuki.NewBot(
		akizuki.DefaultCatalogParser,
		detector,
		akizuki.DefaultParseItem(time.Second),
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	toot, err := akizuki.DefaultToot(config.Mastodon)
	if err != nil {
		log.Fatalln(err.Error())
	}

	shot := akizuki.ScreenShotOneTakeScreenShot(config.ScreenshotOne)

	bot.SetFormat(
		func(item *akizuki.Item) string {
			return fmt.Sprintf("[%v] %v %v\n%v", item.ItemCode, item.Name, item.Price, item.Url)
		},
	).SetToot(toot).SetTakeScreenshot(shot)

	fmt.Println(bot.RunOnce())
}
