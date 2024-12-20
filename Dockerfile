FROM golang:alpine AS builder

WORKDIR /steamsalty

RUN apk update && apk add --no-cache ca-certificates 

COPY . .
RUN go mod download
RUN go build -o /app/steamsalty



FROM scratch
WORKDIR /app

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /usr/share/ca-certificates /usr/share/ca-certificates

COPY --from=builder /app/steamsalty /app/steamsalty

ENTRYPOINT ["/app/steamsalty"]