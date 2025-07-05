# --- Build stage ---
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/app ./cmd/api

# --- Final stage ---
FROM gcr.io/distroless/static-debian12 AS final

COPY --from=builder /app/app /app

ARG PORT=8081
ENV PORT=${PORT}

EXPOSE ${PORT}

WORKDIR /
USER nonroot:nonroot

CMD ["/app"]
