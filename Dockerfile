
FROM golang:1.21-alpine AS builder

ARG VERSION=1.21


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

RUN go get -d -v ./cmd/api \
    && CGO_ENABLED=0 go build -ldflags="-s -w -X 'main.Version=${VERSION}'" -o /app/memnixrest ./cmd/api/main.go\
    && upx /app/memnixrest

RUN wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.5/dumb-init_1.2.5_x86_64
RUN chmod +x /usr/local/bin/dumb-init

FROM gcr.io/distroless/static:nonroot AS production

LABEL org.opencontainers.image.source=https://github.com/memnix/memnix-rest
LABEL description="Production stage for Memnix REST API."


ENV TZ Europe/Paris

WORKDIR /app

COPY --from=builder  /app/memnixrest /app/memnixrest
COPY --from=builder  /build/favicon.ico /app/favicon.ico
COPY --from=builder  /usr/local/bin/dumb-init /usr/bin/dumb-init


EXPOSE 1815

USER nonroot

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/app/memnixrest"]
