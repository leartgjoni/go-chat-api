FROM golang:1.13.10-alpine3.11

ARG PORT
ARG REDIS_ADDR
ENV PORT=${PORT}
ENV REDIS_ADDR=${REDIS_ADDR}

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./build/go-chat-api ./cmd/app

# This container exposes port 8080 to the outside world
EXPOSE ${PORT}

# Run the binary program produced by `go install`
CMD ["./build/go-chat-api"]