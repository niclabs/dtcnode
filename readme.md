# ZMQ Node for  PKCS11-Compatible DTC Library

A Golang rework of [our C++ Library](https://github.com/niclabs/tchsm-libdtc) in Golang, but using a PKCS#11 interface.

This node is used in our implementation of [PKCS11-Compatible DTC Library with ZMQ communication module](https://github.com/niclabs/dtc).


# Installation

## System Requirements (Recommended)

* Debian 10.2 (Buster)
* 1 GB of RAM
* At least 1 GB of available local storage (for keyshare database)

## How to build

The following libraries should be installed in the systems which are going to use the compiled library:

* git
* tar
* wget
* libzmq3-dev v4 or greater (for zmq communication with the nodes)
* libczmq-dev (for zmq communication with the nodes)
* gcc
* Go (1.13.4 or higher)

On [Debian 10 (Buster)](https://www.debian.org), with a sudo-enabled user, the commands to run to install dependencies and 
build are the following:

```bash
# Install requirements
sudo apt install libzmq3-dev libczmq-dev build-essential pkg-config git tar wget
```

Then, you need to install Go 1.13.4 or higher. You can find how to install Go on [its official page](https://golang.org/doc/install).

The following command allows you to clone this repository.
```
# Clone and compile repository
git clone https://github.com/niclabs/dtcnode
```

Finally, you can build the program, executing the following commands:

```
cd dtcnode
go build
```

The program will be named `dtcnode` and will be compiled on the same folder as the cloned git repository. 

# Getting Configuration and Key Pairs

To execute `dtcnode`, you first need the configuration files.

You should use the files generated with [dtcconfig](https://github.com/niclabs/dtcconfig) 
when you built the [dHSM Library](https://github.com/niclabs/dtc) and configured it.
Copy the config corresponding to this node in `/etc/dtcnode` folder.

## Docker Nodes set

The `docker-compose` file on `docker` folder is useful to test the DTC library with nodes deployed on the same machine.
Assuming that `docker` and `docker-compose` are already installed, you need to start the containers with:

`./docker/test.sh`

Then, you should copy the generated `config.yaml` file in `/etc/dtc/config.yaml`. The database file is going to be by default
in `/tmp/dtc.sqlite3`.

When you want to stop the nodes, you can simply stop the `./docker/test.sh` script.