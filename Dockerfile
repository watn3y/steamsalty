FROM golang:alpine AS builder

WORKDIR /tmp/build

COPY . .
RUN go mod download
RUN go build -o /tmp/build/bin/steamsalty



FROM scratch
WORKDIR /app
COPY --from=builder /tmp/build/bin/steamsalty /app/steamsalty


ENTRYPOINT ["/app/steamsalty"]