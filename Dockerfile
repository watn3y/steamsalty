FROM --platform=$BUILDPLATFORM golang:1.25.2-alpine AS builder

WORKDIR /steamsalty

RUN apk update && apk add --no-cache ca-certificates 

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o ./steamsalty



FROM scratch
LABEL org.opencontainers.image.source=https://github.com/watn3y/steamsalty
LABEL org.opencontainers.image.description="SteamSalty"
LABEL org.opencontainers.image.licenses=GPL-3.0
WORKDIR /app

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /usr/share/ca-certificates /usr/share/ca-certificates

COPY --from=builder /steamsalty/steamsalty /app/steamsalty

ENTRYPOINT ["/app/steamsalty"]
