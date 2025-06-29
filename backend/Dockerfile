# --- Build stage ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api

# --- Final stage ---
FROM alpine:3.18

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/app .

RUN chown appuser:appgroup /app/app

USER appuser

EXPOSE 8081

CMD ["./app"]
