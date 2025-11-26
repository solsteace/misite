FROM golang:alpine AS base
WORKDIR /site

# Tools obtained this way: Goose
FROM alpine/curl:latest AS utils-from-curl
WORKDIR .
RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh

FROM oven/bun:alpine AS ts-watcher
RUN mkdir -p /temp/dev
COPY package.json bun.lock tsconfig.json /temp/dev
RUN cd /temp/dev && bun install --frozen-lockfile
CMD ["bun", "run", "dev"]

FROM base AS devel
COPY --from=utils-from-curl /usr/local/bin /
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", "./.air.toml"]