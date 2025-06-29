package scrapper

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

func ExtractStructureResources(e *colly.HTMLElement) (string, int) {
	return ExtractStructureResourceName(e), ExtractStructureResourceQuantity(e)
}

func ExtractStructureResourceName(e *colly.HTMLElement) string {
	resourceName := e.ChildText("a[title]")
	if resourceName == "" {
		slog.Warn("Resource name not found", "html", e.Text)
		return ""
	}

	if resourceName == "Polymer, Organic Polymer, or Corrupted Nodule" {
		resourceName = "Polymer"
	}

	return resourceName
}

func ExtractStructureResourceQuantity(e *colly.HTMLElement) int {
	var quantityStr string

	e.DOM.Contents().EachWithBreak(func(i int, s *goquery.Selection) bool {
		if len(s.Nodes) > 0 && s.Nodes[0].Type == html.TextNode {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				quantityStr = text
				return false
			}
		}
		return true
	})

	if quantityStr == "" {
		slog.Warn("Could not find quantity text", "html", e.Text)
		return 0
	}
	quantityStr = strings.TrimSpace(strings.Replace(quantityStr, "×", "", -1))
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		slog.Error("Failed to parse quantity", "error", err)
		return 0
	}

	return quantity
}
