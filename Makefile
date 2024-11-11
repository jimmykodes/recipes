all: dist dist/style.css

.PHONY: dist
dist:
	go run ./cmd/build build --route-prefix "/recipes"

.PHONY: local
local:
	go run ./cmd/build build --dist ./local
	$(MAKE) local/style.css

dist/style.css local/style.css: assets/style/app.styl
	stylus < $^ > $@

serve:
	python3 -m http.server -d local

clean:
	rm -rf dist local
