FROM golang:1.19-alpine as builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update && apk add tzdata && apk add upx

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .
RUN go get -d -v
RUN go build -ldflags="-s -w" -o /app/memnixrest .
RUN upx /app/memnixrest

FROM alpine:3.17

RUN apk update && apk add ca-certificates
COPY --from=builder /usr/share/zoneinfo/Europe/Paris /usr/share/zoneinfo/Europe/Paris
ENV TZ Europe/Paris

WORKDIR /app

COPY --from=builder /app/memnixrest /app/memnixrest
COPY --from=builder /build/.env /app/.env

EXPOSE 1815

RUN apk add dumb-init
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/app/memnixrest"]
