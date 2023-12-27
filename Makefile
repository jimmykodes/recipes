all:
	go run ./cmd/build build
	mkdir dist/static
	stylus < assets/style/app.styl > dist/static/style.css

prod:
	go run ./cmd/build build --route-prefix "/recipes"
	mkdir dist/static
	stylus < assets/style/app.styl > dist/static/style.css
