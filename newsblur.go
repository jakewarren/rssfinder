package rssfinder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type newsBlurResult []newsBlurResultElement

func unmarshalNewsBlurResult(data []byte) (newsBlurResult, error) {
	var r newsBlurResult
	err := json.Unmarshal(data, &r)
	return r, err
}

type newsBlurResultElement struct {
	NumSubscribers int64  `json:"num_subscribers"`
	Tagline        string `json:"tagline"`
	Value          string `json:"value"`
	Label          string `json:"label"`
	ID             int64  `json:"id"`
}

func NewsBlurSearch(u string) []Feed {
	feeds := make([]Feed, 0)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("https://newsblur.com/rss_feeds/feed_autocomplete?term=%s", u))
	if err != nil {
		log("error searching newsblur: %s", err.Error())
		return feeds
	}

	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	i, unmarshalErr := unmarshalNewsBlurResult(data)
	if unmarshalErr != nil {
		log("error unmarshaling result from newsblur: %s", unmarshalErr.Error())
		return feeds
	}

	for _, item := range i {
		feedURL := item.Value

		feeds = append(feeds, Feed{FeedURL: feedURL, Sources: []string{"newsblur"}})
		log("got feed from newsblur: %s", feedURL)

	}

	return feeds
}
