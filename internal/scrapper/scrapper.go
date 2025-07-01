package scrapper

import (
	"log/slog"
	"strings"

	"github.com/gocolly/colly"
)

func GetStructuresLinks(c *colly.Collector, structuresCollector *colly.Collector) {
	c.OnHTML("tbody tr td a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if strings.HasPrefix(link, "/wiki/File:") {
			return
		}

		structuresCollector.Visit(e.Request.AbsoluteURL(link))
	})
}

func GetStructureDetails(c *colly.Collector, structures []Item) {
	c.OnHTML("div.info-framework", func(e *colly.HTMLElement) {
		resourcesText := e.ChildText("div[style*='padding-left:5px'] b")
		baseResourcesText := e.ChildText("div.mw-collapsible-content tr td")

		pageTitle := e.ChildText("div[class*='info-arkitex info-X1-100 info-masthead']")
		resources := ExtractStructureResources(resourcesText)
		baseResources := ExtractStructureResources(baseResourcesText)

		if len(resources) != 0 && pageTitle != "" {
			structure := Item{
				Name:          pageTitle,
				Resources:     resources,
				BaseResources: baseResources,
			}

			structures = append(structures, structure)
		}
	})
}

func SetupCollectorLogs(c *colly.Collector, name string) {
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "ASA-Resource-Calculator/1.0")
		slog.Info("Visiting", "collector", name, "url", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		slog.Error("Error visiting", "collector", name, "url", r.Request.URL.String(), "error", err)
	})
}
