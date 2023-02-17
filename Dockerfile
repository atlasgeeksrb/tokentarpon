#build
FROM golang:1.20.1-alpine3.17 as builder
WORKDIR /usr/local/go/src/
COPY . ./
WORKDIR /usr/local/go/src/tokenizerService
RUN go build .

# post build
FROM alpine:3.17

RUN mkdir -p /var/tokenapp

WORKDIR /var/tokenapp
COPY --from=builder /usr/local/go/src/tokenizerService .

RUN addgroup app && adduser -S -G app app
RUN chown -R app:app /var/tokenapp
USER app
EXPOSE 8090

CMD ["./tokenizerService"]
