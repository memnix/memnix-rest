FROM golang:1.21-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata && apk add --no-cache upx

WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY default.pgo .
RUN go mod download

COPY . .

RUN go get -d -v \
    && CGO_ENABLED=0 go build -pgo=auto -ldflags="-s -w" -o /app/memnixrest .\
    && upx /app/memnixrest

RUN wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.5/dumb-init_1.2.5_x86_64
RUN chmod +x /usr/local/bin/dumb-init

FROM gcr.io/distroless/static:nonroot AS production

COPY --from=builder  /usr/share/zoneinfo/Europe/Paris /usr/share/zoneinfo/Europe/Paris
ENV TZ Europe/Paris

ENV GOMEMLIMIT 4000MiB

WORKDIR /app

COPY --from=builder  /app/memnixrest /app/memnixrest
COPY --from=builder  /build/.env* /app/.
COPY --from=builder  /build/favicon.ico /app/favicon.ico
COPY --from=builder  /build/config /app/config
COPY --from=builder  /usr/local/bin/dumb-init /usr/bin/dumb-init


EXPOSE 1815

USER nonroot

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/app/memnixrest"]
