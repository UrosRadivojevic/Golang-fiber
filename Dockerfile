# Use an official Go image as the base image
FROM golang:1.19-alpine

# Set the working directory in the container
WORKDIR /app

# Copy the source code to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Specify the command to run when the container starts
CMD ["./main"]

# RUN go get github.com/urosradivojevic/health/handler
# RUN go get github.com/urosradivojevic/health/container