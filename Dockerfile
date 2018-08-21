FROM golang:1.9-alpine as builder
WORKDIR src/git.containerum.net/ch/mail-templater
COPY . .
WORKDIR cmd/mail-templater
RUN CGO_ENABLED=0 go build -v -tags "jsoniter" -ldflags="-w -s -extldflags '-static'" -o /bin/mail-templater
COPY templates.db /tmp/templates.db

FROM alpine:3.7
RUN apk --no-cache add ca-certificates

VOLUME ["/storage"]

# app
COPY --from=builder /bin/mail-templater /
COPY --from=builder /tmp/templates.db /storage/

# timezone data
ENV GIN_MODE=debug \
    CH_MAIL_LOG_LEVEL=4 \
    CH_MAIL_TEMPLATE_DB="/storage/templates.db" \
    CH_MAIL_MESSAGES_DB="/storage/messages.db" \
    CH_MAIL_UPSTREAM=smtp \
    CH_MAIL_UPSTREAM_SIMPLE=smtp \
    CH_MAIL_SENDER_NAME_SIMPLE=containerum \
    CH_MAIL_SENDER_MAIL_SIMPLE=noreply-test@containerum.io \
    CH_MAIL_SENDER_NAME=containerum \
    CH_MAIL_SENDER_MAIL=noreply-test@containerum.io \
    CH_MAIL_USER_MANAGER_URL=http://user-manager:8111 \
    CH_MAIL_LISTEN_ADDR=:7070 \
    CH_MAIL_SMTP_ADDR=mail:465 \
    CH_MAIL_SMTP_LOGIN=noreply-test@containerum.io \
    CH_MAIL_SMTP_PASSWORD=verystrongpassword

EXPOSE 7070

ENTRYPOINT ["/mail-templater"]
