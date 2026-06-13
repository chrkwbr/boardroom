# syntax=docker/dockerfile:1.7

FROM golang:1.25 AS builder
WORKDIR /workspace
COPY . .

WORKDIR /workspace/chat/cmd/api-query
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/api-query .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/api-query /app
EXPOSE 8081
ENTRYPOINT ["/app"]
