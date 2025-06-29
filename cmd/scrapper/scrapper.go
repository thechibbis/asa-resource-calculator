package main

import (
	"log/slog"

	"github.com/gocolly/colly"
)

type Structure struct {
	Name     string
	Resource map[string]int
}

func main() {
	// var structures []Structure

	c := colly.NewCollector(
		colly.AllowedDomains("ark.wiki.gg"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "ASA-Resource-Calculator/1.0")
		slog.Info("Visiting", "url", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		slog.Error("Error visiting", "url", r.Request.URL.String(), "error", err)
	})

	c.OnHTML("span.mw-page-title-main", func(e *colly.HTMLElement) {
		slog.Info("Structure", "name", e.Text)
	})

	c.OnHTML("div.info-X1-100 div[style*='padding-left:5px'] b", func(e *colly.HTMLElement) {
		slog.Info("Resource", "name", e.Text)
	})

	c.Visit("https://ark.wiki.gg/wiki/Small_Tek_Teleporter")
	// c.Visit("https://ark.wiki.gg/wiki/Structures")
}
