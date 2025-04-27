FROM golang:1.24-alpine3.20 AS builder

WORKDIR /app

# comment this command if your service do not build in china
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

# comment this command if your service do not build in china
RUN sed 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' -i /etc/apk/repositories

RUN apk add git

RUN --mount=type=cache,target=/go,id=go go install github.com/gone-io/gonectr@latest && \
    go install go.uber.org/mock/mockgen@latest

COPY ["go.mod", "go.sum", "./"]

RUN --mount=type=cache,target=/go,id=go go mod download

COPY cmd cmd
COPY internal internal

#编译
RUN --mount=type=cache,target=/go,id=go  go generate ./... && \
    go build -ldflags="-w -s" -tags musl -o bin/server ./cmd/server

FROM alpine:3.20

# comment this command if you don't need to change the timezone
RUN sed 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' -i /etc/apk/repositories && \
    apk update && \
    apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
COPY config config
COPY --from=builder /app/bin/server /app/server

ARG ENVIRONMENT=dev
ENV ENV=${ENVIRONMENT}

CMD ["./server"]
EXPOSE 8080