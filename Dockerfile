# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS build

WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags "-s -w" -o /out/fx-rate-api ./cmd/api


FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=build /out/fx-rate-api /app/fx-rate-api
COPY --from=build /src/migrations /app/migrations

ENV HTTP_ADDR=":8080"
ENV MIGRATIONS_DIR="/app/migrations"

EXPOSE 8080

ENTRYPOINT ["/app/fx-rate-api"]
