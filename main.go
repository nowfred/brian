package main

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

var (
	folder = "data"
)

func encode(url string) string {
	uEnc := base64.URLEncoding.EncodeToString([]byte(url))
	return uEnc
}

func collect(domain string) {
	c := colly.NewCollector(
		colly.AllowedDomains(domain),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		uri := e.Request.AbsoluteURL(link)
		c.Visit(uri)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		t := time.Now().UTC().UnixNano()
		h := encode(r.Request.URL.String())
		f := fmt.Sprintf("%s/%d_%s.html", folder, t, h)
		r.Save(f)
	})

	c.Visit(domain)
}

func main() {

	domains := []string{"deadspin.com", "nakedcapitalism.com", "arstechnica.com"}

	for _, d := range domains {
		go collect(d)
	}
	select {}
}
