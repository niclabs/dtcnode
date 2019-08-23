# ZMQ Node for  PKCS11-Compatible DTC Library

A Golang rework of [our C++ Library](https://github.com/niclabs/tchsm-libdtc) in Golang, but using a PKCS#11 interface.

This node is used in our implementation of [PKCS11-Compatible DTC Library with ZMQ communication module](https://github.com/niclabs/dtc).


# Installation

1. Make sure you have ZMQ 4.0 and CZMQ 4.0 or greater installed on your node machine. Also, you will need to have Golang 1.12 or greater to compile the project. 

1. Other packages you need to install on your system are `pkgconfig`, `gcc` and `musl-dev`. They are used in cgo compilation (ZMQ requires them). You can see a config example for Debian Buster in the `Dockerfile` of [DTC `integration_test` folder](https://github.com/niclabs/dtc).
1. Clone this repository.
1. Execute `go mod tidy` to download the dependencies of this project.
1. Build the project executing `go build` in the root of the project. This will create a `dtcnode` executable.
1. If you need a keypair for your server, you can use `dtcnode generate-curve` command to create it.
1. execute `dtcnode generate-config` to create a config file. For more information about how to use this command, check at the end of this readme.
1. Copy the configuration to the current directory, or to `/etc/dtcnode/config.yaml`.
1. Launch the node executing `./dtcnode serve`.

# Creating Configuration and Key Pairs

As mentioned in installation, the `dtcnode` utillity includes two useful commands related to the node confguration:

## `dtcnode generate-config`

`genconfig` generates a dtcnode configuration. 

It is used with the following arguments: 
 1. `n` as the node IP and listening port, with an `:` between both values. _eg: 192.168.0.20:2030_
 1. `c` as the client IP. _eg: 192.168.0.4_
 1. `k` as the server public key in Base85 encoding format
 1. `o` as the output location for the config file (by default is the current working directory).

**Example** `dtcnode generate-config -n 192.168.0.22:2030 -c 192.168.0.22 -k {0j3IXL0Jw:)K$b1@(1=<8z/joPM.c+EXVBMS>7$ -o ./config.yaml`

You can get more information executing `dtcnode generate-config help`.

## `dtcnode generate-curve`

 Prints to stdout a random public and private key usable on a ZMQ server or node.

It has no arguments.

**Example** `dtcnode generate-curve`


## Docker Tests

The `docker-compose` file on `docker` folder is useful to test the DTC library with nodes deployed on the same machine. Assuming that `docker` and `docker-compose` are already installed, you need to start the containers with:

```bash
docker-compose build
docker-compose up
```
