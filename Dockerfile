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


ENV DB_HOST=localhost
ENV DB_PORT=3306
ENV DB_USER=root
ENV DB_PASSWORD=password
ENV DB_NAME=jwt_demo
ENV JWT_SECRET=invisiblekey!

WORKDIR /app/

COPY --from=0 /app/auth_service ./

EXPOSE 8080

# Command to run the executable
CMD ["./auth_service"]


