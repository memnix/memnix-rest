
FROM golang:1.21-alpine3.19 AS builder

ARG VERSION=prod

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk --no-cache add tzdata upx

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w -X 'main.Version=${VERSION}'" -o /app/memnixrest ./cmd/api/main.go \
    && upx /app/memnixrest

RUN wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.5/dumb-init_1.2.5_x86_64
RUN chmod +x /usr/local/bin/dumb-init

RUN apk del upx

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
