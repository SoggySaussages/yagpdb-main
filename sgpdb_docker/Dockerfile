FROM golang:1.20.6 as builder

WORKDIR /appbuild/sgpdb
COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /appbuild/sgpdb/cmd/sgpdb
RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags "-X github.com/botlabs-gg/sgpdb/v2/common.VERSION=$(git describe --tags)"



FROM alpine:latest

WORKDIR /app
VOLUME ["/app/soundboard", "/app/cert"]
EXPOSE 80 443

# Dependencies: ca-certificates for client TLS, tzdata for timezone and ffmpeg for soundboard support
RUN apk --no-cache add ca-certificates ffmpeg tzdata

COPY --from=builder /appbuild/sgpdb/cmd/sgpdb/sgpdb sgpdb

ENTRYPOINT ["/app/sgpdb"]
CMD ["-all", "-pa", "-exthttps=false", "-https=true"]
