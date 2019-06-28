# ZMQ Node for  PKCS11-Compatible DTC Library

A Golang rework of [our C++ Library](https://github.com/niclabs/tchsm-libdtc) in Golang, but using a PKCS#11 interface.

This node is used in our implementation of [PKCS11-Compatible DTC Library with ZMQ communication module](https://github.com/niclabs/dtc).


# Installation

1. Make sure you have ZMQ 4.0 and CZMQ 4.0 or greater installed on your node machine.
1. Other packages you need to install on your system are `pkgconfig`, `gcc` and `musl-dev`. They are used in cgo compilation (ZMQ requires them). You can see a config example in Alpine Linux in the `Dockerfile` of [DTC `integration_test` folder](https://github.com/niclabs/dtc).
1. Clone this repository.
1. Build the project executing `go build` in the root of the project. This will create a `dtcnode` executable.
1. If you need a keypair for your server, you can use `gencurve` executable (build it executing `go build` in `gencurve` folder) to create a new keypair.
1. Build the go project in the folder `genconfig` and execute `genconfig` with the following arguments: 
 1. `n` as the node IP and listening port, with an `:` between both values. _eg: 192.168.0.20:2030_
 1. `s`as the server IP and listening port, with an `:`between both values. _eg: 192.168.0.4:2030_
 1. `k`as the server public key in Base85 encoding format.
1. Copy the configuration to the current directory, or to `/etc/dtc/config.yaml`.
1. Launch the node executing `./dtcnode`.

# Creating Configuration and Key Pairs

As mentioned in installation, the source code includes two useful utilities related to the nodes:

- genconfig: generates a dtcnode configuration, based on servers information.
- gencurve: generates a pair of CURVE keys used by ZMQ to secure the communication channel.
