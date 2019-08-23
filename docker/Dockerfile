FROM golang:buster

MAINTAINER Eduardo Riveros<eduardo@niclabs.cl>

ENV DTC_NODE 0.0.0.0:2030
ENV DTC_CLIENT 0.0.0.0:3030
ENV DTC_PK changemeoriwontwork

WORKDIR /
RUN apt update
RUN apt -y install libzmq3-dev libczmq-dev build-essential musl-dev

# copy again again
COPY ./dtcnode/ /dtcnode

WORKDIR /dtcnode
RUN go mod tidy -v
RUN go build

RUN mkdir /etc/dtcnode

CMD /dtcnode/dtcnode generate-config -n $DTC_NODE -c $DTC_CLIENT -k $DTC_PK -o /etc/dtcnode/config.yaml && /dtcnode/dtcnode serve
