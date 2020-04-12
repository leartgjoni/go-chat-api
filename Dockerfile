FROM golang:1.13.10-alpine3.11 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./build/go-chat-api ./cmd/app


FROM alpine

COPY --from=builder /app/build/go-chat-api /app/go-chat-api

ARG PORT
ARG REDIS_ADDR
ENV PORT=${PORT}
ENV REDIS_ADDR=${REDIS_ADDR}

# This container exposes port ${PORT} to the outside world
EXPOSE ${PORT}

# Run the binary program produced by `go install`
ENTRYPOINT ["/app/go-chat-api"]