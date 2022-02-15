package rssfinder

// LogFunc is a function that logs the provided message with optional
// fmt.Sprintf-style arguments. By default, discards logs.
var LogFunc func(string, ...interface{}) = nil

type RSSFinder struct {
	Config struct { // allows the enabling or disabling of "active" searching methods
		EnableFuzzer    bool // fuzzes for common feed URLs
		EnableScraper   bool // scrapes the URL and homepage for RSS feeds
		EnableFeedly    bool // queries Feedly's API for RSS feed URLs
		EnableInoreader bool // queries Inoreader's API for RSS feed URLs
		EnableNewsBlur  bool // queries NewsBlur's API for RSS feed URLs
		EnableSmartMode bool // will run the scraper and fuzzer mode if the three APIs don't return anything
	}
}

type Feed struct {
	FeedURL string
	Sources []string
}

func New() *RSSFinder {
	return &RSSFinder{
		Config: struct {
			EnableFuzzer    bool
			EnableScraper   bool
			EnableFeedly    bool
			EnableInoreader bool
			EnableNewsBlur  bool
			EnableSmartMode bool
		}{EnableFuzzer: true, EnableScraper: false, EnableFeedly: true, EnableInoreader: true, EnableNewsBlur: true, EnableSmartMode: true},
	}
}

func log(msg string, v ...interface{}) {
	if LogFunc != nil {
		LogFunc(msg, v...)
	}
}

func (r *RSSFinder) FindRSS(url string) []Feed {
	feeds := make([]Feed, 0)

	feeds = append(feeds, r.scrapeURL(url)...)
	if r.Config.EnableInoreader {
		feeds = append(feeds, InoreaderSearch(url)...)
	}
	if r.Config.EnableFeedly {
		feeds = append(feeds, FeedlySearch(url)...)
	}
	if r.Config.EnableNewsBlur {
		feeds = append(feeds, NewsBlurSearch(url)...)
	}
	if r.Config.EnableFuzzer {
		feeds = append(feeds, fuzzURL(url)...)
	}

	return unique(feeds)
}

// unique the URLs but keep information about all sources which indicate that URL
func unique(orig []Feed) []Feed {
	encountered := map[string]Feed{}

	for _, i := range orig {
		if _, ok := encountered[i.FeedURL]; ok {
			for range encountered[i.FeedURL].Sources {
				newSources := appendIfMissing(encountered[i.FeedURL].Sources, i.Sources[0])
				encountered[i.FeedURL] = Feed{
					FeedURL: i.FeedURL,
					Sources: newSources,
				}
			}
		} else {
			// record the element as encountered
			encountered[i.FeedURL] = i
		}
	}

	result := make([]Feed, 0)

	for k := range encountered {
		result = append(result, encountered[k])
	}

	return result
}

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
