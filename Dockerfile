# Stage 1: Build the Go application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .
# Install tzdata to get time zone info
RUN apt-get update && apt-get install -y tzdata

# Download dependencies and build the Go app
RUN go mod tidy
RUN go build -o app .

# Stage 2: Create the final image for the Go application
FROM alpine:3.20


# Set the working directory inside the container
WORKDIR /app

# Install tzdata to get time zone info
RUN apt-get update && apt-get install -y tzdata

# Copy the built Go application and the .env file from the builder stage
COPY --from=builder /app/app /app
COPY .env /app/.env


# Expose the port the application will run on
EXPOSE 80

# Command to run the Go application
CMD ["./app"]
