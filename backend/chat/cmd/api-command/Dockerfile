# syntax=docker/dockerfile:1.7

FROM golang:1.25 AS builder
WORKDIR /workspace
COPY . .

WORKDIR /workspace/chat/cmd/api-command
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/api-command .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/api-command /app
EXPOSE 8080
ENTRYPOINT ["/app"]
