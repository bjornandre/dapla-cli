

FROM golang:alpine AS builder

WORKDIR /build

RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/statisticsnorway/dapla-cli/

COPY main.go .
COPY cmd ./cmd
COPY rest ./rest

RUN go get -d -v

ENV CGO_ENABLED=0 \
#    GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

RUN go build -o /dapla .

FROM scratch

COPY --from=builder /dapla /dapla
#CMD ["/dapla"]
ENTRYPOINT ["/dapla"]