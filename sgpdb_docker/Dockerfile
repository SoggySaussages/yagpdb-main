FROM docker.io/golang:1.23.3 AS builder

WORKDIR /appbuild/sgpdb
COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /appbuild/sgpdb/cmd/sgpdb
RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags "-X github.com/SoggySaussages/sgpdb/common.VERSION=$(git describe --tags)"



FROM docker.io/alpine:latest

WORKDIR /app
VOLUME ["/app/soundboard", "/app/cert"]
EXPOSE 80 443

# Dependencies: ca-certificates for client TLS, tzdata for timezone and ffmpeg for soundboard support
RUN apk --no-cache add ca-certificates ffmpeg tzdata

COPY --from=builder /appbuild/sgpdb/cmd/sgpdb/sgpdb sgpdb

ENTRYPOINT ["/app/sgpdb"]
CMD ["-all", "-pa", "-exthttps=false", "-https=true"]
