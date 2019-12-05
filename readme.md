# ZMQ Node for  PKCS11-Compatible DTC Library

A Golang rework of [our C++ Library](https://github.com/niclabs/tchsm-libdtc) in Golang, but using a PKCS#11 interface.

This node is used in our implementation of [PKCS11-Compatible DTC Library with ZMQ communication module](https://github.com/niclabs/dtc).


# Installation

# How to build

First, it's necessary to download all the requirements of the Go project. The following libraries should be installed in the systems which are going to use the library:

* libzmq v4 or greater (for zmq communication with the nodes)
* libczmq (for zmq communication with the nodes)
* gcc
* Go (1.13.4 or higher)

for building the project as a library, you should execute the following command. It will produce a file named `dtcnode` that you can execute.

`go build`

On Ubuntu 18.04 LTS, the commands to run to build are the following:

```bash
# Install requirements
sudo apt install libzmq3-dev libczmq-dev build-essential pkg-config

# Download and install Go 1.13.4 (or higher) for Linux AMD 64 bit.
wget https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.13.4.linux-amd64.tar.gz

# ADD /usr/local/go to PATH
export PATH=$PATH:/usr/local/go

# Clone and compile repository
git clone https://github.com/niclabs/dtcnode
cd dtcnode
./build.sh
```

# Getting Configuration and Key Pairs

You should use the files generated with [dtcconfig](https://github.com/niclabs/dtcconfig) 
when you built the [dHSM Library](https://github.com/niclabs/dtc). Copy the config corresponding to this node in
/etc/dtcnode folder.

## Docker Tests

The `docker-compose` file on `docker` folder is useful to test the DTC library with nodes deployed on the same machine. Assuming that `docker` and `docker-compose` are already installed, you need to start the containers with:

`./docker/test.sh`
