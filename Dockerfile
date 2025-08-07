FROM --platform=$BUILDPLATFORM golang:1.24.6-alpine AS builder

WORKDIR /steamsalty

RUN apk update && apk add --no-cache ca-certificates 

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o ./steamsalty



FROM scratch
WORKDIR /app

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /usr/share/ca-certificates /usr/share/ca-certificates

COPY --from=builder /steamsalty/steamsalty /app/steamsalty

ENTRYPOINT ["/app/steamsalty"]
