FROM golang:buster

MAINTAINER Eduardo Riveros<eduardo@niclabs.cl>

WORKDIR /
RUN apt update
RUN apt -y install libzmq3-dev libczmq-dev build-essential musl-dev

# copy again again
COPY ./dtcnode/ /dtcnode

WORKDIR /dtcnode
RUN go mod tidy -v
RUN go build

RUN mkdir /etc/dtcnode

CMD /dtcnode/dtcnode
