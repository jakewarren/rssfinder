package rssfinder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type inoreaderResult []inoreaderResultElement

func unmarshalInoreaderResult(data []byte) (inoreaderResult, error) {
	var r inoreaderResult
	err := json.Unmarshal(data, &r)
	return r, err
}

type inoreaderResultElement struct {
	Value    *string `json:"value,omitempty"`
	Type     *string `json:"type,omitempty"`
	Label    *string `json:"label,omitempty"`
	ItemType *string `json:"item_type,omitempty"`
	ID       *int64  `json:"id,omitempty"`
}

func InoreaderSearch(u string) []Feed {
	feeds := make([]Feed, 0)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("https://www.inoreader.com/autocomplete.php?term=%s&origin=smart_search", u))
	if err != nil {
		log("error searching inoreader: %s", err.Error())
		return feeds
	}

	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	i, unmarshalErr := unmarshalInoreaderResult(data)
	if unmarshalErr != nil {
		log("error unmarshaling result from inoreader: %s", unmarshalErr.Error())
		return feeds
	}

	for _, item := range i {
		if item.Type != nil {
			if *item.Type == "feed" {
				feeds = append(feeds, Feed{FeedURL: *item.Value, Sources: []string{"inoreader"}})
				log("got feed from inoreader: %s", *item.Value)
			}
		}
	}

	return feeds
}
