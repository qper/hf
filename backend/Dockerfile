FROM golang:1.25.11-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY cmd/server ./cmd/server
RUN go build -o /out/server ./cmd/server

FROM alpine:3.20
COPY --from=builder /out/server /server
ENTRYPOINT ["/server"]
