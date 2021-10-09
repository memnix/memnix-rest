FROM golang:latest

RUN mkdir -p /go/src/memnix-rest
WORKDIR /go/src/memnix-rest

COPY . /go/src/memnix-rest

RUN go get -d -v
RUN go install -v

EXPOSE 1813

CMD ["/go/bin/memnixrest"]
