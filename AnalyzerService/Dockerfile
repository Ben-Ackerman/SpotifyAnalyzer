# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.13.0-alpine as builder

ENV GO111MODULE=on

# Add Maintainer Info
LABEL maintainer="Benjamin Ackerman <ben@tacmail.net>"

# Set the Current Working Directory inside the container
WORKDIR /SpotifyAnalyzer/app

# Copy go mod and sum files
COPY AnalyzerService/go.mod AnalyzerService/go.sum ./

COPY . /SpotifyAnalyzer/

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY AnalyzerService/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /SpotifyAnalyzer/app/main .

COPY AnalyzerService/src/ ./src

# Command to run the executable
ENTRYPOINT ["./main"]