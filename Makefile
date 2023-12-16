all:
	go run ./cmd/build -dist ./docs
	mkdir docs/static
	stylus < assets/style/app.styl > docs/static/style.css
