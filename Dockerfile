FROM golang:1.26-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/url-shortener ./cmd/url-shortener

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /bin/url-shortener /app/url-shortener
COPY config/ /app/config/

EXPOSE 8082

USER nonroot:nonroot

ENTRYPOINT ["/app/url-shortener"]
