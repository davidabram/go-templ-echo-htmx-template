.PHONY: run templ watch tailwind setup setup-templ

watch:
	@./bin/air & $(MAKE) tailwind

tailwind:
	@./node_modules/.bin/tailwindcss -i ./style.css -o ./dist/style.css --watch

templ:
	@./bin/templ generate ./internals/templates/*.templ

setup-templ:
	@go install github.com/a-h/templ && cp $(shell go env GOPATH)/bin/templ ./bin

.ONESHELL:
setup:
	@curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
	$(MAKE) setup-templ
	@bun i
