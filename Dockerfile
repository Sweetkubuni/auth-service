FROM golang:1.19-alpine as builder

# Create app directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source code into container
COPY . .


# Build th Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o auth_service ./.

FROM alpine:latest

WORKDIR /app/

COPY --from=0 /app/auth_service ./

# Command to run the executable
CMD ["./auth_service"]


