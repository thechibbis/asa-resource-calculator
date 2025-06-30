package scrapper

import (
	"log/slog"
	"regexp"
	"strconv"
	"strings"
)

func ExtractStructureResources(resourcesText string) map[string]int {
	resources := make(map[string]int)

	pattern := `(\d+)\s*×\s*([^\d]*)`
	re, err := regexp.Compile(pattern)
	if err != nil {
		slog.Error("regex error", "re", err)
	}

	matches := re.FindAllStringSubmatch(resourcesText, -1)

	for _, match := range matches {
		resourceName := strings.TrimSpace(match[2])
		quantityStr := match[1]

		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			slog.Error("Failed to convert to a quantity int", "err", err)
			continue
		}

		resources[resourceName] = quantity
	}

	return resources
}
