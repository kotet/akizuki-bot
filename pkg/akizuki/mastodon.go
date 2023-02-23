package akizuki

import (
	"context"
	"io"
	"log"

	"github.com/mattn/go-mastodon"
)

type DefaultTootConfig struct {
	Server       string `json:"server"`
	ClientID     string `json:"id"`
	ClientSecret string `json:"secret"`
	Username     string `json:"mail"`
	Password     string `json:"pass"`
}

func DefaultToot(config DefaultTootConfig) (func(text string, images []io.Reader) error, error) {
	client := mastodon.NewClient(&mastodon.Config{
		Server:       config.Server,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
	})
	err := client.Authenticate(context.Background(), config.Username, config.Password)
	if err != nil {
		return nil, err
	}
	return func(text string, images []io.Reader) error {
		ids := []mastodon.ID{}
		for _, image := range images {
			att, err := client.UploadMediaFromReader(context.Background(), image)
			if err != nil {
				return err
			}
			log.Printf("Image Uploaded: %v %v", att.ID, att.URL)
			ids = append(ids, att.ID)
		}
		s, err := client.PostStatus(context.Background(), &mastodon.Toot{
			Status:     text,
			MediaIDs:   ids,
			Visibility: mastodon.VisibilityUnlisted,
			Language:   "ja",
		})
		if err != nil {
			return err
		}
		log.Printf("Toot Completed: %v", s.URL)
		return nil
	}, nil
}
