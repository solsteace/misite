FROM golang:alpine AS base
WORKDIR /site

# Tools obtained this way: Goose
FROM alpine/curl:latest AS utils-obtained-with-curl
WORKDIR .
RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh

FROM base AS devel
COPY --from=utils-obtained-with-curl /usr/local/bin /
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", "./.air.toml"]