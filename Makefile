.PHONY: dev dev/server dev/templ dev/css generate build 

dev:
	$(MAKE) -j3 dev/css  dev/templ dev/server

dev/server:
	air

dev/templ:
	templ generate --watch

dev/css:
	tailwindcss -i ./web/static/css/input.css -o ./web/static/css/output.css --minify --watch

generate:
	templ generate
	sqlc generate
	tailwindcss -i ./web/static/css/input.css -o ./web/static/css/output.css --minify

build:
	go build -o ./bin/itinera ./cmd/itinera
