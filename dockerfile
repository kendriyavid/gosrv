# # Use the official Go image as a base image
# FROM golang:1.23.5 AS builder

# # Set the working directory inside the container
# WORKDIR /app

# # Copy the Go module files and download dependencies
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy the source code
# COPY . .

# # Build the application
# RUN go build -o gosrv main.go

# # Use a minimal base image for the final container
# FROM alpine:latest

# # Set working directory
# WORKDIR /root/

# # Copy the compiled binary from the builder stage
# COPY --from=builder /app/gosrv .

# # Copy the .env file
# COPY .env .

# # Expose the application port
# EXPOSE 3000

# # Run the application
# CMD ["./gosrv"]


# Use Go official image as a builder
# FROM golang:1.23 as builder

# WORKDIR /app
# COPY . .

# # Download dependencies and build
# RUN go mod tidy
# RUN go build -o main .

# # Use a lightweight image for the final container
# FROM alpine:latest

# WORKDIR /root/
# COPY --from=builder /app/main .

# # Run the application
# CMD ["./main"]

# FROM golang:1.23 as builder

# COPY . .

# # Download dependencies and build
# RUN go mod tidy

# # Build the application
# RUN go build -o gosrv main.go

# # Use a minimal base image for the final container
# FROM alpine:latest

# # Set working directory
# WORKDIR /root/

# # Copy the compiled binary from the builder stage
# COPY --from=builder /app/gosrv .

# # Copy the .env file
# COPY .env .

# # Expose the application port
# EXPOSE 3000

# # Run the application
# CMD ["./gosrv"]

#---------------------------------------------------------------------


# # Build Stage
# FROM golang:1.23.0 AS builder  

# # Set working directory inside the container
# WORKDIR /app

# # Copy go.mod and go.sum first to leverage Docker cache
# COPY go.mod go.sum ./

# # Download dependencies
# RUN go mod download

# # Copy the source code
# COPY . .

# # Build the application
# RUN go build -o gosrv main.go

# # Final Stage (Using Alpine for minimal size)
# FROM alpine:latest


# # Set working directory
# WORKDIR /root/

# # Copy the compiled binary from the builder stage
# COPY --from=builder /app/gosrv .

# # Copy the .env file
# COPY .env . 

# # Expose the application port
# EXPOSE 3000

# # Run the application
# CMD ["./gosrv"]


# Use a specific Go version as the base image
FROM golang:1.23.0

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o main main.go

# Expose the port the app will run on
EXPOSE 3000

# Define the command to start your app
CMD ["./main"]
