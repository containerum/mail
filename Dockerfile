FROM golang:1.11-alpine as builder
RUN apk add --update make git
WORKDIR src/git.containerum.net/ch/mail-templater
COPY . .
RUN VERSION=$(git describe --abbrev=0 --tags) make build-for-docker

FROM alpine:3.8
RUN apk --no-cache add ca-certificates

VOLUME ["/storage"]

# app
COPY --from=builder /tmp/mail /
COPY templates.json /

# timezone data
ENV GIN_MODE=debug \
    TEMPLATE_DB="/storage/templates.db" \
    MESSAGES_DB="/storage/messages.db" \
    DEFAULT_TEMPLATES="templates.json" \
    UPSTREAM=smtp \
    UPSTREAM_SIMPLE=smtp \
    SENDER_NAME_SIMPLE=containerum \
    SENDER_MAIL_SIMPLE=noreply-test@containerum.io \
    SENDER_NAME=containerum \
    SENDER_MAIL=noreply-test@containerum.io \
    USER_MANAGER_URL=http://user-manager:8111 \
    LISTEN_ADDR=:7070 \
    SMTP_ADDR=mail:465 \
    SMTP_LOGIN=noreply-test@containerum.io \
    SMTP_PASSWORD=verystrongpassword

EXPOSE 7070

ENTRYPOINT ["/mail"]
