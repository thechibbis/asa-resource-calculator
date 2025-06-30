package main

import (
	"log/slog"
	"strings"

	"github.com/gocolly/colly"
	"github.com/thechibbis/asa-resource-calculator/internal/scrapper"
)

type Structure struct {
	Name          string
	Resources     map[string]int
	BaseResources map[string]int
}

func main() {
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.AllowedDomains("ark.wiki.gg"),
	)

	details := c.Clone()

	structures := make([]Structure, 0, 200)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "ASA-Resource-Calculator/1.0")
		slog.Info("Visiting", "url", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		slog.Error("Error visiting", "url", r.Request.URL.String(), "error", err)
	})

	details.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "ASA-Resource-Calculator/1.0")
		slog.Info("Details Visiting", "url", r.URL.String())
	})

	details.OnError(func(r *colly.Response, err error) {
		slog.Error("Error visiting", "url", r.Request.URL.String(), "error", err)
	})

	c.OnHTML("tbody tr td a[href]", func(e *colly.HTMLElement) {
		// If attribute class is this long string return from callback
		// As this a is irrelevant
		link := e.Attr("href")

		if strings.HasPrefix(link, "/wiki/File:") {
			return
		}

		// start scaping the page under the link found
		details.Visit(e.Request.AbsoluteURL(link))
	})

	details.OnHTML("div.info-framework", func(e *colly.HTMLElement) {
		resourcesText := e.ChildText("div[style*='padding-left:5px'] b")
		baseResourcesText := e.ChildText("div.mw-collapsible-content tr td")

		pageTitle := e.ChildText("div[class*='info-arkitex info-X1-100 info-masthead']")
		resources := scrapper.ExtractStructureResources(resourcesText)
		baseResources := scrapper.ExtractStructureResources(baseResourcesText)

		if len(resources) != 0 && pageTitle != "" {
			structure := Structure{
				pageTitle,
				resources,
				baseResources,
			}

			structures = append(structures, structure)
		}
	})

	c.Visit("https://ark.wiki.gg/wiki/Structures")
	//details.Visit("https://ark.wiki.gg/wiki/Spinning_Mule_(Primitive_Plus)")
	c.Wait()

	slog.Info("Structures", "st", structures)
}
