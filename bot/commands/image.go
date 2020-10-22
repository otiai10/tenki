package commands

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/otiai10/amesh/bot/slack"
	"github.com/otiai10/goapis/google"
)

// ImageCommand ...
type ImageCommand struct{}

// Match ...
func (cmd ImageCommand) Match(payload *slack.Payload) bool {
	if len(payload.Ext.Words) == 0 {
		return false
	}
	return payload.Ext.Words[0] == "img" || payload.Ext.Words[0] == "image"
}

// Handle ...
func (cmd ImageCommand) Handle(ctx context.Context, payload *slack.Payload) slack.Message {
	client := google.Client{
		APIKey:               os.Getenv("GOOGLE_CUSTOMSEARCH_API_KEY"),
		CustomSearchEngineID: os.Getenv("GOOGLE_CUSTOMSEARCH_ENGINE_ID"),
	}
	words := payload.Ext.Words[1:]
	if len(words) == 0 {
		return slack.Message{Channel: payload.Event.Channel, Text: "検索語が無いです"}
	}
	query := strings.Join(words, "+")
	rand.Seed(time.Now().Unix())
	res, err := client.SearchImage(query, 1+rand.Intn(10))
	if err != nil {
		return slack.Message{Channel: payload.Event.Channel, Text: fmt.Sprintf("%v\n> %v", err.Error(), words)}
	}
	if len(res.Items) == 0 {
		return slack.Message{Channel: payload.Event.Channel, Text: fmt.Sprintf("結果はゼロでした\n> %v", words)}
	}
	// TODO: ランダムにひとつ選ぶ
	item := res.Items[0]
	return slack.Message{
		Channel: payload.Event.Channel,
		Blocks:  []slack.Block{{Type: "image", ImageURL: item.Link, AltText: query}},
	}
}

// Help ...
func (cmd ImageCommand) Help(payload *slack.Payload) slack.Message {
	return slack.Message{
		Channel: payload.Event.Channel,
		Text:    "画像検索コマンド",
	}
}