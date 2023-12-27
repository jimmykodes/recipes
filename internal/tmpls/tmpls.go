package tmpls

import (
	"html/template"
	"path"
	"path/filepath"
	"strings"

	"github.com/jimmykodes/gommand"
	"github.com/jimmykodes/gommand/flags"
	"github.com/jimmykodes/strman"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	templateDirFlag = "template-dir"
	routePrefixFlag = "route-prefix"
)

func Flags() *flags.FlagSet {
	fs := flags.NewFlagSet()

	fs.String(templateDirFlag, "templates", "directory containing template files")
	fs.String(routePrefixFlag, "", "prefix string for generated routes")

	return fs
}

func New(ctx *gommand.Context) (*Tmpls, error) {
	tmpls := Tmpls{
		location:    ctx.Flags().String(templateDirFlag),
		routePrefix: ctx.Flags().String(routePrefixFlag),
	}

	var err error
	tmpls.Template, err = template.New("templates").
		Funcs(tmpls.funcs()).
		ParseGlob(filepath.Join(tmpls.location, "*"))

	return &tmpls, err
}

type Tmpls struct {
	*template.Template
	location    string
	routePrefix string
}

func Encode(s string) string {
	return strman.ToSnake(strings.Join(strings.Fields(s), "_"))
}

func (t Tmpls) route(s string) string {
	return path.Join(t.routePrefix, s)
}

func (t Tmpls) funcs() template.FuncMap {
	caser := cases.Title(language.English)
	return template.FuncMap{
		"title": func(s string) string {
			return caser.String(strings.ReplaceAll(s, "_", " "))
		},
		"route":  t.route,
		"encode": Encode,
	}
}
