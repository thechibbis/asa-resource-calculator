package scrapper

type Item struct {
	Name          string         `json:"name"`
	Resources     map[string]int `json:"resources"`
	BaseResources map[string]int `json:"base_resources"`
}
