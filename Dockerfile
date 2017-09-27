# Docker image for the Drone rsync plugin
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-rsync .

FROM golang:1.8 as builder

WORKDIR /go/src/github.com/russ-p/drone-rsync/

RUN go get github.com/codegangsta/cli &&   \
go get github.com/joho/godotenv/autoload && \
go get golang.org/x/crypto/ssh

ADD main.go        .
ADD plugin.go      .
ADD main_test.go   .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w"


FROM alpine:3.6
RUN apk add -U ca-certificates openssh-client rsync && rm -rf /var/cache/apk/*
COPY  --from=builder /go/src/github.com/russ-p/drone-rsync/drone-rsync /bin/
ENTRYPOINT ["/bin/drone-rsync"]
