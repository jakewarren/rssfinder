package rssfinder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type feedlyResult struct {
	Results   []feedlyMetadata `json:"results"`
	QueryType string           `json:"queryType"`
	Scheme    string           `json:"scheme"`
}

type feedlyMetadata struct {
	FeedID      string `json:"feedId"`
	LastUpdated int64  `json:"lastUpdated"`
	ID          string `json:"id"`
	Title       string `json:"title"`
	Updated     int64  `json:"updated"`
	Subscribers int64  `json:"subscribers"`
	Website     string `json:"website"`
}

func unmarshalFeedlyResult(data []byte) (feedlyResult, error) {
	var r feedlyResult
	err := json.Unmarshal(data, &r)
	return r, err
}

func FeedlySearch(u string) []Feed {
	feeds := make([]Feed, 0)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("https://feedly.com/v3/search/feeds?q=%s", u))
	if err != nil {
		log("error searching feedly: %s", err.Error())
		return feeds
	}

	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	i, unmarshalErr := unmarshalFeedlyResult(data)
	if unmarshalErr != nil {
		log("error unmarshaling result from feedly: %s", unmarshalErr.Error())
		return feeds
	}

	for _, item := range i.Results {
		feedURL := strings.TrimPrefix(item.FeedID, "feed/")

		feeds = append(feeds, Feed{FeedURL: feedURL, Sources: []string{"feedly"}})
		log("got feed from feedly: %s", feedURL)

	}

	return feeds
}
