FROM golang:buster

MAINTAINER Eduardo Riveros<eduardo@niclabs.cl>

ENV DTC_NODE 0.0.0.0:2030
ENV DTC_SERVER 0.0.0.0:3030
ENV DTC_PK changemeoriwillnotwork

WORKDIR /
RUN apt update
RUN apt -y install libzmq3-dev libczmq-dev build-essential musl-dev

RUN git clone https://github.com/niclabs/dtcnode --branch reqsocket

WORKDIR /dtcnode
RUN go mod tidy -v
RUN go build

WORKDIR /dtcnode/genconfig
RUN go mod tidy -v
RUN go build

RUN mkdir /etc/dtcnode

CMD /dtcnode/genconfig/genconfig -n $DTC_NODE -s $DTC_SERVER -k $DTC_PK -o /etc/dtcnode/config.yaml && /dtcnode/dtcnode
