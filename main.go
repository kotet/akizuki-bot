package main

import (
	"encoding/json"
	"fmt"
	"log"

	_ "embed"

	"github.com/kotet/akizuki-bot/pkg/akizuki"
)

//go:embed config.json
var configJSON []byte

type Config struct {
	Mastodon      akizuki.DefaultTootConfig   `json:"mastodon"`
	ScreenshotOne akizuki.ScreenShotOneConfig `json:"screenshotone"`
}

func main() {

	var config Config
	err := json.Unmarshal(configJSON, &config)
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
