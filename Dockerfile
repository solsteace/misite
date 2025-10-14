FROM golang:alpine AS base
WORKDIR /site

FROM base AS devel
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", "./.air.toml"]