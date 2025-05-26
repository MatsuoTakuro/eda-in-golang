FROM golang:1.24-alpine AS builder

WORKDIR /mallbots

RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY .. ./

# Build the binary.
RUN go build -gcflags="all=-N -l" -o monolith ./cmd/mallbots

FROM alpine:latest AS runtime

RUN apk add --no-cache libc6-compat
COPY --from=builder /go/bin/dlv /go/bin/dlv

COPY --from=builder /mallbots/docker/wait-for .
RUN chmod +x /wait-for

COPY --from=builder /mallbots/docker/entrypoint .
RUN chmod +x /entrypoint

# Copy the binary to the production image from the builder stage.
COPY --from=builder /mallbots/monolith /mallbots/monolith
