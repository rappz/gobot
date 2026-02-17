FROM golang:1.24.2-bookworm
# Set the working directory
WORKDIR /app
# Copy the go.mod and go.sum files
COPY . .
# Download the dependencies
RUN go mod download
# Builds the app
RUN go build -o /chrisbot
#Execute the binary
CMD ["/chrisbot"]