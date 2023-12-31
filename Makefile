.PHONY: run template watch setup air

watch:
	@./bin/air & $(MAKE) tailwind

tailwind:
	@./node_modules/.bin/tailwindcss -i ./style.css -o ./dist/style.css --watch

templ:
	@./bin/templ generate ./internals/templates/*.templ

tailwind-build:
	@./node_modules/.bin/tailwindcss -i ./style.css -o ./dist/style.css

.ONESHELL:
setup:
	@curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
	@go install github.com/a-h/templ/cmd/templ@latest && cp $(shell go env GOPATH)/bin/templ ./bin
	@bun i
