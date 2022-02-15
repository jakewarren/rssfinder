package rssfinder

import (
	"net/url"
	"strings"
)

var suffixes = []string{
	// Generic suffixes
	"/index.xml", "/atom.xml", "/feed.xml", "/feeds", "/feeds/default", "/feed", "/feed/", "/feed/default",
	"/feeds/posts/default/", "?feed=rss", "?feed=atom", "?feed=rss2", "?feed=rdf", "/rss",
	"/atom", "/rdf", "/index.rss", "/index.rdf", "/index.atom",
	"?type=100",                       // Typo3 RSS URL
	"/blog-feed.xml",                  // Wix.com RSS URL
	"?format=rss",                     // squarespace.com format
	"?format=feed&type=rss",           // Joomla RSS URL
	"/feeds/posts/default",            // Blogger.com RSS URL
	"/rss.xml",                        // Posterous.com RSS feed
	"/articles.rss", "/articles.atom", // Patch.com RSS feeds
	"/commits/master.atom", "/commits/main.atom", // github RSS feeds
	"/.rss", // reddit
}

func fuzzURL(u string) []Feed {
	feeds := make([]Feed, 0)

	if !strings.HasPrefix(u, "http") {
		u = "http://" + u
	}
	parsedURL, _ := url.Parse(u)

	// fuzz the url the user provided
	feeds = append(feeds, fuzz(u)...)

	// if the user specified a url, scrape the root domain as well
	if u != parsedURL.Scheme+parsedURL.Host {
		feeds = append(feeds, fuzz(parsedURL.Scheme+"://"+parsedURL.Host)...)
	}

	return feeds
}

// fuzz the url to try to find any feeds.
//    this is purposely single-threaded to be more polite to websites. probably not needed but..
func fuzz(u string) []Feed {
	feeds := make([]Feed, 0)

	for _, s := range suffixes {
		var fuzzedURL string

		// ensure we don't get double slashes in the path
		if strings.HasSuffix(u, "/") && strings.HasPrefix(s, "/") {
			fuzzedURL = u + strings.TrimPrefix(s, "/")
		} else {
			fuzzedURL = u + s
		}

		if !strings.HasPrefix(fuzzedURL, "http") {
			fuzzedURL = "http://" + fuzzedURL
		}

		f := checkLink(fuzzedURL)
		if f != "" {
			feeds = append(feeds, Feed{FeedURL: f, Sources: []string{"fuzzer"}})
			log("[fuzzer] found feed %s", f)
		}
	}

	return feeds
}
