all: dist dist/style.css

local: local/*.html local/style.css

.PHONY: dist
dist:
	go run ./cmd/build build --route-prefix "/recipes"

local/*.html: recipes/*.yaml
	go run ./cmd/build build --dist ./local

dist/style.css local/style.css: assets/style/app.styl
	stylus < $^ > $@

serve:
	python3 -m http.server -d local

clean:
	rm -rf dist local
