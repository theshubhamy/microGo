# Use Go 1.24-alpine for the build stage
FROM golang:1.24-alpine AS build

# Install required dependencies
RUN apk --no-cache add gcc g++ make ca-certificates

# Set the working directory in the container to the Go source directory
WORKDIR /go/src/github.com/theshubhamy/microGo

# Copy go.mod and go.sum for Go modules
COPY go.mod go.sum ./

# Copy the vendor folder (if using vendoring)
COPY vendor vendor

# Copy the entire services directory, which includes catalog, account, order, etc.
COPY services services

# Build the order service application
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./services/order/cmd

# Final stage: use a minimal Alpine image
FROM alpine:3.11

# Set the working directory
WORKDIR /usr/bin

# Copy the built app from the build stage
COPY --from=build /go/bin .

# Expose the port the app will run on
EXPOSE 8080

# Run the app
CMD ["app"]
