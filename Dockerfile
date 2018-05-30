FROM golang
MAINTAINER Sylvain Laurent

ENV GOBIN $GOPATH/bin
ENV PROJECT_DIR github.com/Magicking/faktur-daemon
ENV PROJECT_NAME faktur

ADD vendor /usr/local/go/src
ADD cmd /go/src/${PROJECT_DIR}/cmd
ADD merkle /go/src/${PROJECT_DIR}/merkle
ADD common /go/src/${PROJECT_DIR}/common
ADD backends /go/src/${PROJECT_DIR}/backends
ADD internal /go/src/${PROJECT_DIR}/internal

WORKDIR /go/src/${PROJECT_DIR}

RUN go build -v -o /go/bin/main /go/src/${PROJECT_DIR}/cmd/${PROJECT_NAME}-server/main.go
ENTRYPOINT /go/bin/main
