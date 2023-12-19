package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jimmykodes/strman"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var (
	templateDir = flag.String("templates", "./templates", "directory containing template files")
	recipeDir   = flag.String("recipes", "./recipes", "directory containing recipe files")
	distDir     = flag.String("dist", "./dist", "directory to write html files to")
)

type Recipe struct {
	Title       string    `yaml:"title"`
	Tags        []string  `yaml:"tags"`
	Ingredients []Section `yaml:"ingredients"`
	Procedure   []Section `yaml:"procedure"`
}

type Section struct {
	Name  string   `yaml:"name"`
	Items []string `yaml:"items"`
}

type Link struct {
	Title string
	Ref   string
}

type LinkPage struct {
	Title string
	Links []Link
}

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func encode(s string) string {
	return strman.ToSnake(strings.Join(strings.Fields(s), "_"))
}

func run() error {
	flag.Parse()

	caser := cases.Title(language.English)
	fm := template.FuncMap{
		"title": func(s string) string {
			return caser.String(strings.ReplaceAll(s, "_", " "))
		},
		"encode": encode,
	}

	tmpl, err := template.New("templates").Funcs(fm).ParseGlob(filepath.Join(*templateDir, "*.go.html"))
	if err != nil {
		return err
	}

	if err := os.RemoveAll(*distDir); err != nil {
		return err
	}
	if err := os.Mkdir(*distDir, 0755); err != nil {
		return err
	}

	var links []Link
	tags := make(map[string][]Link)

	err = filepath.WalkDir(*recipeDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".yaml" {
			return nil
		}

		r, err := os.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()

		var recipe Recipe
		if err := yaml.NewDecoder(r).Decode(&recipe); err != nil {
			return err
		}

		ext := filepath.Ext(path)
		fn := filepath.Base(path)
		fn = strings.TrimSuffix(fn, ext)
		fn = strman.ToSnake(fn) + ".html"

		f, err := os.Create(filepath.Join(*distDir, fn))
		if err != nil {
			return err
		}
		defer f.Close()

		link := Link{Title: recipe.Title, Ref: fn}
		links = append(links, link)

		for _, tag := range recipe.Tags {
			t := encode(tag)
			tags[t] = append(tags[t], link)
		}

		return tmpl.ExecuteTemplate(f, "recipe.go.html", recipe)
	})
	for tag, links := range tags {
		err := func(tag string, links []Link) error {
			f, err := os.Create(filepath.Join(*distDir, "tag_"+tag+".html"))
			if err != nil {
				return err
			}
			defer f.Close()
			return tmpl.ExecuteTemplate(f, "tag.go.html", LinkPage{Title: tag, Links: links})
		}(tag, links)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(*distDir, "index.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	slices.SortFunc(links, func(a, b Link) int {
		switch {
		case a.Title < b.Title:
			return -1
		case a.Title > b.Title:
			return 1
		default:
			return 0
		}
	})

	return tmpl.ExecuteTemplate(f, "index.go.html", struct {
		Title string
		Links []Link
	}{Title: "Recipes", Links: links})
}
