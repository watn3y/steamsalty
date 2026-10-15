FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

WORKDIR /steamsalty

RUN apk update && apk add --no-cache ca-certificates 

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o ./steamsalty



FROM scratch

LABEL org.opencontainers.image.title="SteamSalty" \
    org.opencontainers.image.description="SteamSalty notifies you on telegram about new comments on any steam profile, with built in auto translation." \
    org.opencontainers.image.authors="watn3y <me@watn3y.de>" \
    org.opencontainers.image.vendor="watn3y" \
    org.opencontainers.image.documentation="https://git.watn3y.de/watn3y/steamsalty/src/branch/master/README.md" \
    org.opencontainers.image.url="https://git.watn3y.de/watn3y/steamsalty" \
    org.opencontainers.image.source="https://git.watn3y.de/watn3y/steamsalty" \
    org.opencontainers.image.licenses="GPL-3.0-or-later"

WORKDIR /app

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /usr/share/ca-certificates /usr/share/ca-certificates

COPY --from=builder /steamsalty/steamsalty /app/steamsalty

ENTRYPOINT ["/app/steamsalty"]
