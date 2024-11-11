package recipes

type Recipe struct {
	Title       string    `yaml:"title"`
	Category    []string  `yaml:"category"`
	Tags        []string  `yaml:"tags"`
	Description string    `yaml:"description"`
	Yield       string    `yaml:"yield"`
	Source      string    `yaml:"source"`
	Ingredients []Section `yaml:"ingredients"`
	Procedure   []Section `yaml:"procedure"`
}

type Section struct {
	Name  string   `yaml:"name"`
	Items []string `yaml:"items"`
}
