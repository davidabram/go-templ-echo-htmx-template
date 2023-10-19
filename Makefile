.PHONY: run template watch setup air

run:
	@go run cmd/web/main.go

.ONESHELL:
watch:
	@./bin/air

tailwind:
	@./node_modules/.bin/tailwindcss -i ./style.css -o ./dist/style.css --watch

templ:
	@./bin/templ generate ./internals/templates/*.templ

.ONESHELL:
setup:
	@curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
	@go install github.com/a-h/templ && cp $(shell go env GOPATH)/bin/templ ./bin
	@bun i
