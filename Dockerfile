# Use the official Golang image as a parent image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the Go application code to the container
COPY . .
RUN go build -v -o /usr/local/bin/app ./...

# Expose the port that your Gin application will run on
EXPOSE 8080

# Command to run the application
CMD ["app"]
