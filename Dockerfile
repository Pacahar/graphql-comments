FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o graphql-comments ./cmd/graphql-comments




FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/graphql-comments .
COPY --from=builder /app/config/config.yaml ./config/config.yaml

EXPOSE 4000

CMD ["./graphql-comments"]