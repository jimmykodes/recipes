package cmds

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jimmykodes/gommand"
	"github.com/jimmykodes/gommand/flags"
	"github.com/jimmykodes/recipes/internal/dist"
	"github.com/jimmykodes/recipes/internal/recipes"
	"github.com/jimmykodes/recipes/internal/tmpls"
	"github.com/jimmykodes/strman"
	"gopkg.in/yaml.v3"
)

type Link struct {
	Title string
	Ref   string
}

type LinkPage struct {
	Title string
	Links []Link
}

var buildCmd = &gommand.Command{
	Name: "build",
	FlagSet: flags.NewFlagSet().
		AddFlagSet(dist.Flags()).
		AddFlagSet(tmpls.Flags()),
	Run: func(ctx *gommand.Context) error {
		distDir, err := dist.New(ctx)
		if err != nil {
			return err
		}
		if err := distDir.Clear(); err != nil {
			return err
		}

		tmpl, err := tmpls.New(ctx)
		if err != nil {
			return err
		}

		var links []Link
		tags := make(map[string][]Link)

		if err := filepath.WalkDir("./recipes", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".yaml" {
				return nil
			}

			r, err := os.Open(path)
			if err != nil {
				return err
			}
			defer r.Close()

			var recipe recipes.Recipe
			if err := yaml.NewDecoder(r).Decode(&recipe); err != nil {
				return err
			}

			ext := filepath.Ext(path)
			fn := filepath.Base(path)
			fn = strings.TrimSuffix(fn, ext)
			fn = strman.ToSnake(fn) + ".html"

			f, err := os.Create(filepath.Join(distDir.Dir(), fn))
			if err != nil {
				return err
			}
			defer f.Close()

			link := Link{Title: recipe.Title, Ref: fn}
			links = append(links, link)

			for _, tag := range recipe.Tags {
				t := tmpls.Encode(tag)
				tags[t] = append(tags[t], link)
			}

			return tmpl.ExecuteTemplate(f, "recipe.go.html", recipe)
		}); err != nil {
			return err
		}

		for tag, links := range tags {
			err := func(tag string, links []Link) error {
				f, err := distDir.Create("tag_" + tag + ".html")
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

		tagLinks := make([]Link, 0, len(tags))
		for tag := range tags {
			tagLinks = append(tagLinks, Link{Title: tag, Ref: "tag_" + tag + ".html"})
		}
		slices.SortFunc(tagLinks, func(a, b Link) int {
			switch {
			case a.Title < b.Title:
				return -1
			case a.Title > b.Title:
				return 1
			default:
				return 0
			}
		})

		tagFile, err := distDir.Create("tags.html")
		if err != nil {
			return err
		}
		defer tagFile.Close()

		if err := tmpl.ExecuteTemplate(tagFile, "tags.go.html", LinkPage{Title: "Tags", Links: tagLinks}); err != nil {
			return err
		}

		indexFile, err := distDir.Create("index.html")
		if err != nil {
			return err
		}
		defer indexFile.Close()

		if err := tmpl.ExecuteTemplate(indexFile, "index.go.html", LinkPage{Title: "Recipes", Links: links}); err != nil {
			return err
		}
		return nil
	},
}
