FROM oven/bun:1.0.7-debian as bun

WORKDIR /app

COPY package.json bun.lockb tailwind.config.js tsconfig.json style.css ./internals/templates/ ./

RUN bun install

RUN bunx tailwindcss -i ./style.css -o ./dist/style.css

FROM golang:1.21.3-bullseye AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

RUN useradd -u 1001 crocoder

COPY . .

RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o web \
  ./cmd/web/main.go

FROM scratch

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=bun /app/dist/style.css /dist/style.css

COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /app/web /web

COPY --from=build /app/.env /.env

USER crocoder

EXPOSE 3000

CMD ["/web"]

