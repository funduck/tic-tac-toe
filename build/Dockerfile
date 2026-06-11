# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS build_base
WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY ./cmd cmd
COPY ./internal internal
COPY ./docs docs

RUN swag init --parseDependency -g cmd/server/main.go

RUN go build  -o ./server ./cmd/server/main.go

FROM alpine:3.18
WORKDIR /app

# Create non-root user and switch to it
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Copy built binary with user permissions
COPY --from=build_base --chown=appuser:appuser /app/server ./server
COPY --from=build_base --chown=appuser:appuser /app/docs ./docs

CMD ["./server"]