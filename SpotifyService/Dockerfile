# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.13.0-alpine as builder

ENV GO111MODULE=on

# Add Maintainer Info
LABEL maintainer="Benjamin Ackerman <ben@tacmail.net>"

# Set the Current Working Directory inside the container
WORKDIR /SpotifyAnalyzer/app

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apk add --no-cache ca-certificates

# Copy go mod and sum files
COPY SpotifyService/go.mod SpotifyService/go.sum ./

COPY . /SpotifyAnalyzer/

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY SpotifyService/ .

# Build the Go app
# Compile the Go application as a static binary by setting CGO_ENABLED=0 
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-s -w' -a -o main .


######## Start a new stage from scratch #######
FROM scratch  

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /SpotifyAnalyzer/app/main .

COPY SpotifyService/src/ ./src

# Command to run the executable
ENTRYPOINT ["./main"]