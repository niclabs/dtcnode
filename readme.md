# ZMQ Node for  PKCS11-Compatible DTC Library

A Golang rework of [our C++ Library](https://github.com/niclabs/tchsm-libdtc) in Golang, but using a PKCS#11 interface.

This node is used in our implementation of [PKCS11-Compatible DTC Library with ZMQ communication module](https://github.com/niclabs/dtc).


# Installation

# Creating Configuration and Key Pairs

This source code includes two useful utilities related to the nodes:

- genconfig: generates a dtcnode configuration, based on servers information.
- gencurve: generates a pair of CURVE keys used by ZMQ to secure the communication channel.
