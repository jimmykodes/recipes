all: static
	go run ./cmd/build build

prod: static
	go run ./cmd/build build --route-prefix "/recipes"


.PHONY: static
static: dist/static/style.css

dist/static/style.css: assets/style/app.styl
	mkdir dist/static
	stylus < assets/style/app.styl > dist/static/style.css

serve:
	python3 -m http.server -d dist
