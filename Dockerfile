# Build stage
FROM golang:alpine AS build-env
# Set the working directory for copying the source code
WORKDIR /go/src/pelucio
# Copy the source code into the container's working directory
COPY . .
# Build the Go application
RUN go build -o app

# Final stage
FROM alpine
# Copy the compiled binary from the build stage to the final stage
COPY --from=build-env /go/src/pelucio/app /usr/bin/pelucio
# Set the entry point for the container
ENTRYPOINT ["pelucio"]
