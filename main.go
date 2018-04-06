package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gocolly/colly"
)

var (
	folder = "data"
)

func encode(url string) string {
	uEnc := base64.URLEncoding.EncodeToString([]byte(url))
	return uEnc
}

// Site contains all needed information to start scraping a domain
// and login if necessary
type Site struct {
	Site           string
	AllowedDomains []string
	Initial        string
	Hostname       string
	LoginPage      string
	Username       string
	Password       string
	FormUser       string
	FormPass       string
	ExtraFormData  []string
}

// Config holds all
type Config struct {
	Site []Site
}

func auth(site Site, c *colly.Collector) (*colly.Collector, error) {
	// No-op if no login
	if site.LoginPage == "" {
		fmt.Println("No auth needed for", site.Site)
		return c, nil
	}
	// Assemble post data
	d := map[string]string{}
	d[site.FormUser] = site.Username
	d[site.FormPass] = site.Password
	for _, pair := range site.ExtraFormData {
		kv := strings.Split(pair, ":")
		d[kv[0]] = d[kv[1]]
	}

	err := c.Post(site.LoginPage, d)
	if err != nil {
		return c, err
	}
	return c, nil
}

func collect(site Site) {
	c := colly.NewCollector(
		colly.AllowedDomains(site.AllowedDomains...),
	)

	// Login if we can
	var err error
	c, err = auth(site, c)
	if err != nil {
		fmt.Println("Auth failure for", site.Site)
		return
	}

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
		if r.StatusCode == 200 {
			t := time.Now().UTC().UnixNano()
			h := encode(r.Request.URL.String())
			f := fmt.Sprintf("%s/%d_%s.html", folder, t, h)
			r.Save(f)
		} else {
			fmt.Println("ERROR ON", r.Request.URL.String())
		}
	})
	c.Visit(site.Initial)
}

func main() {
	// Load toml configuration
	var config Config
	if _, err := toml.DecodeFile("conf.toml", &config); err != nil {
		log.Fatal(err)
	}

	// Create collector for each site in config
	for _, s := range config.Site {
		go collect(s)
	}
	select {}
}
