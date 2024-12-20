FROM golang:alpine AS builder

WORKDIR /tmp/build

RUN apk update && apk add --no-cache ca-certificates 

COPY . .
RUN go mod download
RUN go build -o /tmp/build/bin/steamsalty



FROM scratch
WORKDIR /app
COPY --from=builder /tmp/build/bin/steamsalty /app/steamsalty
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /usr/share/ca-certificates /usr/share/ca-certificates

ENTRYPOINT ["/app/steamsalty"]