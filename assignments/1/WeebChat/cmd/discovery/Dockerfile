FROM golang:latest AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/WeebChat

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ../../

# Unit tests
RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN make buildDiscoveryService

# Start fresh from a smaller image
FROM alpine:3.9 
RUN apk add ca-certificates

COPY --from=build_base /tmp/WeebChat/bin/discoveryService /app/discoveryService

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
ENV HOST=localhost
ENV PORT=8080
CMD ["/app/discoveryService"]