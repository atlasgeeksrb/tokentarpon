# build tokenizer Go service
FROM golang:1.20.1-alpine3.17 AS build

WORKDIR /usr/local/go/src/
COPY . ./
WORKDIR /usr/local/go/src/tokenizerService
RUN go build .

# build app image using alpine and the tokenizer app
FROM alpine:3.17

RUN addgroup app && adduser -S -G app app
RUN mkdir -p /var/tokenapp
RUN mkdir -p /var/tokenapp/tokenizerService

WORKDIR /var/tokenapp/tokenizerService
COPY --from=build /usr/local/go/src/tokenizerService .

WORKDIR /var/tokenapp
COPY --from=build /usr/local/go/src/config.json ./config.json

RUN chown -R app:app /var/tokenapp
USER app
EXPOSE 8090

CMD ["./tokenizerService/tokenizerService"]

