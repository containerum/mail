FROM golang:1.9-alpine as builder
WORKDIR src/git.containerum.net/ch/user-manager
COPY . .
RUN CGO_ENABLED=0 go build -v -tags "jsoniter" -ldflags="-w -s -extldflags '-static'" -o /bin/user-manager

FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

FROM scratch
# app
COPY --from=builder /bin/user-manager /
# migrations
COPY migrations /migrations
# timezone data
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# tls certificates
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV GIN_MODE=release \
    CH_USER_LOG_LEVEL=4 \
    CH_USER_PG_URL="postgres://user:password@postgres:5432/user_manager?sslmode=disable" \
    CH_USER_MAIL_URL="http://ch-mail-templater:8080/" \
    CH_USER_RECAPTCHA_KEY="recaptcha_key" \
    CH_USER_GITHUB_APP_ID="github_app" \
    CH_USER_GITHUB_SECRET="github_secret" \
    CH_USER_GOOGLE_APP_ID="google_app" \
    CH_USER_GOOGLE_SECRET="google_secret" \
    CH_USER_FACEBOOK_APP_ID="facebook_app" \
    CH_USER_FACEBOOK_SECRET="facebook_secret" \
    CH_USER_AUTH_GRPC_ADDR="ch-auth:8000" \
    CH_USER_WEB_API_URL="http://web-api:8080" \
    CH_USER_LISTEN_ADDR=":8111"
ENTRYPOINT ["/user-manager"]
