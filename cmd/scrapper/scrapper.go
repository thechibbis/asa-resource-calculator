package main

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/gocolly/colly"
	"github.com/thechibbis/asa-resource-calculator/internal/scrapper"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("ark.wiki.gg"),
		colly.CacheDir("./cache"),
	)

	scrapper.SetupCollectorLogs(c, "MainCollector")

	structures := make([]scrapper.Item, 0)

	structureCollector := c.Clone()
	structureDetailsCollector := c.Clone()

	scrapper.SetupCollectorLogs(structureCollector, "Structures")
	scrapper.SetupCollectorLogs(structureDetailsCollector, "StructureDetails")

	structureCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "ASA-Resource-Calculator/1.0")
		slog.Info("Details Visiting", "url", r.URL.String())
	})

	structureCollector.OnError(func(r *colly.Response, err error) {
		slog.Error("Error visiting", "url", r.Request.URL.String(), "error", err)
	})

	scrapper.GetStructuresLinks(structureCollector, structureDetailsCollector)
	scrapper.GetStructureDetails(structureDetailsCollector, structures)

	c.Visit("https://ark.wiki.gg/wiki/Structures")

	c.Wait()

	data, err := json.Marshal(structures)
	if err != nil {
		slog.Error("Couldn't convert the items to JSON", "err", err)
		return
	}

	os.WriteFile("items.json", data, 0644)
}
