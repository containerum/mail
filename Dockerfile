FROM golang:1.9-alpine as builder
WORKDIR src/bitbucket.org/exonch/ch-auth
COPY . .
RUN CGO_ENABLED=0 go build -v -tags "jsoniter" -ldflags="-w -s -extldflags '-static'" -o /bin/ch-mail-templater

FROM scratch
# app
COPY --from=builder /bin/ch-mail-templater /
# timezone data
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# tls certificates
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV GIN_MODE=release \
    CH_MAIL_LOG_LEVEL=4 \
    CH_MAIL_TEMPLATE_DB=/storage/template.db \
    CH_MAIL_MESSAGES_DB=/storage/messages.db \
    CH_MAIL_UPSTREAM=mailgun \
    MG_API_KEY=apikey \
    MG_DOMAIN=domain \
    MG_PUBLIC_API_KEY=pubkey \
    MG_URL=url \
    CH_MAIL_LISTEN_ADDR=:7070
VOLUME ["/storage"]
ENTRYPOINT ["/ch-mail-templater"]
