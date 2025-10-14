FROM golang:alpine AS base
WORKDIR /site

FROM alpine/curl:latest AS devtools
WORKDIR .
RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh

FROM base AS devel
COPY --from=devtools /usr/local/bin /
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", "./.air.toml"]