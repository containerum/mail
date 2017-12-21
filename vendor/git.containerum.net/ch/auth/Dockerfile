FROM golang:1.9-alpine as builder
WORKDIR src/git.containerum.net/ch/auth
COPY . .
RUN CGO_ENABLED=0 go build -v -ldflags="-w -s -extldflags '-static'" -o /bin/auth

FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

FROM scratch
# app
COPY --from=builder /bin/auth /
# timezone data
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# tls certificates
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV CH_AUTH_HTTP_LISTENADDR=0.0.0.0:8080 \
    CH_AUTH_GRPC_LISTENADDR=0.0.0.0:8888 \
    CH_AUTH_LOG_MODE=text \
    CH_AUTH_LOG_LEVEL=4 \
    CH_AUTH_TOKENS=jwt \
    CH_AUTH_JWT_SIGNING_METHOD=HS256 \
    CH_AUTH_ISSUER=containerum.com \
    CH_AUTH_ACCESS_TOKEN_LIFETIME=15m \
    CH_AUTH_REFRESH_TOKEN_LIFETIME=48h \
    CH_AUTH_JWT_SIGNING_KEY_FILE=/keys/jwt.key \
    CH_AUTH_JWT_VALIDATION_KEY_FILE=/keys/jwt.key \
    CH_AUTH_STORAGE=buntdb \
    CH_AUTH_BUNT_STORAGE_FILE=/storage/storage.db \
    CH_AUTH_TRACER=zipkin \
    CH_AUTH_ZIPKIN_COLLECTOR=nop
VOLUME ["/keys", "/storage"]
EXPOSE 8080 8888
ENTRYPOINT ["/auth"]
