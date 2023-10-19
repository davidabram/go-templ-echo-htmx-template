.PHONY: run template

run:
	@go run cmd/web/main.go

template:
	@./bin/templ generate ./internals/templates/*.templ

.ONESHELL:
setup:
	@curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
	@go install github.com/a-h/templ && cp $(shell go env GOPATH)/bin/templ ./bin

