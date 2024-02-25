all: static dist/*.html

dist/*.html: recipes/*.yaml
	go run ./cmd/build build

prod: static
	go run ./cmd/build build --route-prefix "/recipes"


.PHONY: static
static: dist/style.css

dist/style.css: assets/style/app.styl
	stylus < assets/style/app.styl > dist/style.css

serve:
	python3 -m http.server -d dist
