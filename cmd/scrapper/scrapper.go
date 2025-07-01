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

	structures := make([]scrapper.Item, 0)

	structureCollector := c.Clone()
	structureDetailsCollector := c.Clone()

	scrapper.SetupCollectorLogs(c, "MainCollector")
	scrapper.SetupCollectorLogs(structureCollector, "Structures")
	scrapper.SetupCollectorLogs(structureDetailsCollector, "StructureDetails")

	scrapper.GetStructuresLinks(structureCollector, structureDetailsCollector)
	scrapper.GetStructureDetails(structureDetailsCollector, &structures)

	structureCollector.Visit("https://ark.wiki.gg/wiki/Structures")

	structureCollector.Wait()

	data, err := json.Marshal(structures)
	if err != nil {
		slog.Error("Couldn't convert the items to JSON", "err", err)
		return
	}

	os.WriteFile("items.json", data, 0644)
}
