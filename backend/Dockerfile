FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal
COPY queries ./queries
COPY sqlc.yaml ./
COPY tools.go ./
RUN go build -o /out/server ./cmd/server

FROM alpine:3.20
COPY --from=builder /out/server /server
ENTRYPOINT ["/server"]
