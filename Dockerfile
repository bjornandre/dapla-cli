FROM golang:alpine AS builder

ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

RUN apk update && apk add --no-cache git

WORKDIR /build

COPY .git .
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build --ldflags "-X github.com/statisticsnorway/dapla-cli/cmd.Version=localdev \
    -X github.com/statisticsnorway/dapla-cli/cmd.GitSha1Hash=`git rev-parse --short HEAD` \
    -X github.com/statisticsnorway/dapla-cli/cmd.BuildTime=`date -u +%Y-%m-%dT%H:%M:%S%Z`" \
    -o dapla .

WORKDIR /dist

RUN cp /build/dapla .

EXPOSE 3000
ENTRYPOINT ["/dist/dapla"]
