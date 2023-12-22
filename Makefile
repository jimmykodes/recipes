all:
	go run ./cmd/build
	mkdir dist/static
	stylus < assets/style/app.styl > dist/static/style.css

prod:
	go run -ldflags '-X main.Prefix=/recipes' ./cmd/build
	mkdir dist/static
	stylus < assets/style/app.styl > dist/static/style.css
