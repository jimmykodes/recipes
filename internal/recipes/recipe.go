package recipes

type Recipe struct {
	Title       string    `yaml:"title"`
	Category    []string  `yaml:"category"`
	Tags        []string  `yaml:"tags"`
	Ingredients []Section `yaml:"ingredients"`
	Procedure   []Section `yaml:"procedure"`
}

type Section struct {
	Name  string   `yaml:"name"`
	Items []string `yaml:"items"`
}
