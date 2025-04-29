ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app cmd/api-server/main.go


FROM debian:bookworm

# Copy the binary
COPY --from=builder /run-app /usr/local/bin/

# Copy static files
COPY cmd/api-server/public /usr/local/bin/public
WORKDIR /app

CMD ["run-app"]
