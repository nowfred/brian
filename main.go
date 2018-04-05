package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	// Instantiate default collector
	//disallowedDomainOne, _ := regexp.Compile("twitter.com")
	//disallowedDomainTwo, _ := regexp.Compile("google.com")
	//disallowedDomainThree, _ := regexp.Compile("facebook.com")
	disallowedDomains := []string{
		"twitter.com", "google.com",
		"google.es", "facebook.com",
		"nakedcapitalism.com/author", "apple.com",
		"mozilla.org", "wsj.com", "youtube.com", "pixels.com", "shopify.com",
		"pinterest.com", "footbie.com", "linkedin.com", ".wikimedia.", "wikipedia.", "digg.com", "myspace.com"}

	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.MaxDepth(1),
		colly.AllowedDomains("deadspin.com"),
		//colly.Async(true),
	)

	//c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		uri := e.Request.AbsoluteURL(link)
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		shouldI := true
		for _, d := range disallowedDomains {
			if strings.Contains(uri, d) {
				shouldI = false
			}
		}
		if shouldI {
			c.Visit(uri)
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	//c.Visit("https://nakedcapitalism.com/")
	c.Visit("https://deadspin.com")
	//c.Wait()
}
