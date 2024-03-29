FROM golang:1.21.3-alpine as builder

ENV environment=DEV

WORKDIR /app

COPY ./ /app

RUN go mod download

# Build the binary
RUN go build -o /app/main bff/cmd/server/main.go

# Intermediate stage: Build the binary
FROM golang:1.21.3-alpine

COPY --from=builder /app/main /app/main
COPY --from=builder /app/bff/pkg/config /app/bff/pkg/config
COPY --from=builder /app/docs /app/docs

WORKDIR /app
ENV environment=DEV

EXPOSE 50050

ENTRYPOINT ["./main"]
