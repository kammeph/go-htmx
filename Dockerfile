# ---- Build CSS Stage ----
FROM node:20.11.0 AS cssBuilder
WORKDIR /app
COPY package.json .

RUN npm install
COPY tailwind.config.js .
COPY ./templates ./templates
COPY ./styles/input.css ./styles/input.css
RUN npm run tw:build

# ---- Build Stage ----
FROM golang:1.21.6-alpine AS builder
WORKDIR /app
# Set necessary environmet variables needed for our image
# ENV GO111MODULE=on \
#     CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64
# Move to working directory /build
# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download
# Copy the code into the container
COPY . .
# Build the application
RUN go build -o ./bin/main ./cmd/www

# ---- Production Stage ----
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /app
# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/bin ./bin
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/public ./public
COPY --from=builder /app/i18n ./i18n
COPY --from=cssBuilder /app/styles/global.css ./styles/global.css

# Expose port 8080 to the outside
EXPOSE 8080
# Command to run the executable
CMD ["./bin/main"] 