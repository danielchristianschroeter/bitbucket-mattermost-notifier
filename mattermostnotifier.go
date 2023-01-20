package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type MattermostMessage struct {
	Channel     string       `json:"channel"`
	Username    string       `json:"username"`
	IconEmoji   string       `json:"icon_emoji"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Fallback   string  `json:"fallback"`
	Color      string  `json:"color"`
	Pretext    string  `json:"pretext"`
	Text       string  `json:"text"`
	AuthorName string  `json:"author_name"`
	AuthorIcon string  `json:"author_icon"`
	AuthorLink string  `json:"author_link"`
	Title      string  `json:"title"`
	TitleLink  string  `json:"title_link"`
	Fields     []Field `json:"fields"`
	ImageURL   string  `json:"image_url"`
}

type Field struct {
	Short bool   `json:"short"`
	Title string `json:"title"`
	Value string `json:"value"`
}

func SendMessage(WebhookURL string, channel string, username string, iconEmoji string, text string, attachments []Attachment) error {
	message := MattermostMessage{
		Channel:     channel,
		Username:    username,
		IconEmoji:   iconEmoji,
		Text:        text,
		Attachments: attachments,
	}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", WebhookURL, bytes.NewBuffer(messageJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
