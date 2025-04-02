FROM golang:1.24.2-bookworm
# Set the working directory
WORKDIR /app
# Copy the go.mod and go.sum files
COPY go.mod go.sum main.go bot/bot.go ./
# Download the dependencies
RUN go mod download
# Builds the app
