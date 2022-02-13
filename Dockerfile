FROM golang:1.17-alpine3.15 as goBuilder

COPY ./ /go/src/goddns/

RUN apk update && \
    apk add make build-base && \
    cd /go/src/goddns && \
    go build

FROM node:16-alpine3.15 as nodeBuilder

COPY app /app

RUN cd /app && \
    npm install && \
    npm run build

FROM alpine:3.15

RUN mkdir /goddns

WORKDIR /goddns

COPY --from=goBuilder /go/src/goddns/goddns ./
COPY --from=nodeBuilder /app/dist ./dist

CMD ./goddns