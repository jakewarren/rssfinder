package rssfinder

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var feedContentTypes = []string{
	"application/rss+xml",
	"application/atom+xml",
	"application/rdf+xml",
	"application/rss",
	"application/atom",
	"application/rdf",
	"application/xml",
	"text/rss+xml",
	"text/atom+xml",
	"text/rdf+xml",
	"text/rss",
	"text/atom",
	"text/rdf",
	"text/xml",
}

func scrapeURL(u string) []Feed {
	feeds := make([]Feed, 0)

	if !strings.HasPrefix(u, "http") {
		u = "http://" + u
	}
	parsedURL, _ := url.Parse(u)

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only the specified url, root url, and the root url for any subdomain
		colly.URLFilters(
			regexp.MustCompile(fmt.Sprintf(`^https?://(\w+\.)*?%s/?$`, regexp.QuoteMeta(parsedURL.Host))),
			regexp.MustCompile(fmt.Sprintf(`^%s$`, regexp.QuoteMeta(u))),
		),
		// set the max recursion depth to 1 page
		colly.MaxDepth(1),
		// set a user agent for stealth
		colly.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36`),
	)
	// disable ssl verification
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
	})

	// Scrape all links on page
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// Visit links found on page
		f := checkLink(e.Request.AbsoluteURL(link))
		if f != "" {
			log("[scraper] found feed %s", f)
			feeds = append(feeds, Feed{FeedURL: f, Sources: []string{"scraper"}})
		}

		// Visit links found on page
		// The links are only visited if they are matched by any of the URLFilter regexps
		_ = c.Visit(e.Request.AbsoluteURL(link))
	})

	// Log page visits if verbose output is enabled
	c.OnRequest(func(r *colly.Request) {
		log("Visiting %s", r.URL.String())
	})

	// Start scraping the url
	_ = c.Visit(u)

	// if the user specified a url, scrape the root domain as well
	if u != parsedURL.Scheme+parsedURL.Host {
		_ = c.Visit(parsedURL.Scheme + "://" + parsedURL.Host)
	}

	return feeds
}

var checkedLinks = make(map[string]struct{})

func checkLink(u string) string {
	if !strings.HasPrefix(u, "http") {
		return ""
	}

	// ensure we aren't checking links multiple times
	if _, ok := checkedLinks[u]; ok {
		return ""
	}
	checkedLinks[u] = struct{}{}

	log("checking link for feed content: %s", u)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	req, err := http.NewRequest("HEAD", u, nil)
	if err != nil {
		log("error checking link: %s", err.Error())
		return ""
	}

	req.Header.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36`)

	resp, respErr := client.Do(req)
	if respErr != nil {
		log("error checking link: %s", respErr.Error())
		return ""
	}
	_ = resp.Body.Close()

	if resp.Header == nil {
		return ""
	}

	contentType := resp.Header.Get(`Content-Type`)

	for _, c := range feedContentTypes {
		if strings.Contains(contentType, c) {
			return resp.Request.URL.String()
		}
	}

	return ""
}
