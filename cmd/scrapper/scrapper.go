package main

import (
	"log/slog"

	"github.com/gocolly/colly"
	"github.com/thechibbis/asa-resource-calculator/internal/scrapper"
)

func main() {
	// var structures []Structure

	//_ := make([]scrapper.Structure, 0, 1000)

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

	structure := scrapper.Structure{
		Name:     "",
		Resource: make(map[string]int),
	}

	c.OnHTML("span.mw-page-title-main", func(e *colly.HTMLElement) {
		structure.Name = e.Text
	})

	c.OnHTML("div[style*='padding-left:5px'] > b", func(e *colly.HTMLElement) {
		resourceName, quantity := scrapper.ExtractStructureResources(e)

		if resourceName == "" || quantity <= 0 {
			slog.Warn("Skipping invalid resource", "name", resourceName, "quantity", quantity)
			return
		}

		structure.Resource[resourceName] = quantity
	})

	c.Visit("https://ark.wiki.gg/wiki/Small_Tek_Teleporter")
	// c.Visit("https://ark.wiki.gg/wiki/Structures")

	slog.Info("strucutre", "name", structure.Name, "resources", structure.Resource)
}
