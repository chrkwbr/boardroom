# syntax=docker/dockerfile:1.7

FROM golang:1.25 AS builder
WORKDIR /workspace
COPY . .

WORKDIR /workspace/backend/chat/cmd/consumer-chat
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/consumer-chat .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/consumer-chat /app
ENTRYPOINT ["/app"]
