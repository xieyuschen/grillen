#build stage
FROM golang:alpine AS builder
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache git
RUN apk --no-cache add ca-certificates
WORKDIR /go/src/app
COPY . .
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go get -d -v ./...
RUN apk add build-base 
RUN go build -o /go/bin/app -v main.go
ENTRYPOINT /go/bin/app
