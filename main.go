//go:build embedconfig

package main

import (
	"fmt"
	"log"

	_ "embed"

	"github.com/kotet/akizuki-bot/pkg/akizuki"
	"gopkg.in/yaml.v2"
)

//go:embed config.yaml
var configYAML []byte

type Config struct {
	Mastodon      akizuki.DefaultTootConfig   `yaml:"mastodon"`
	ScreenshotOne akizuki.ScreenShotOneConfig `yaml:"screenshotone"`
}

func main() {

	var config Config
	err := yaml.Unmarshal(configYAML, &config)
	if err != nil {
		log.Fatalln(err.Error())
	}

	bot, err := akizuki.NewBot()
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
