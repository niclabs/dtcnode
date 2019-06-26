package main

import (
	"flag"
	"fmt"
	"github.com/niclabs/dtcnode/config"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
)

var node string
var server string
var pk string
var out string

func init() {
	flag.StringVar(&node, "n", "", "node ip and port")
	flag.StringVar(&server, "s", "", "server ip and port")
	flag.StringVar(&pk, "k", "", "base85 server public key")
	flag.StringVar(&out, "o", "./config.yaml", "path where to output the config file (default: ./config.yaml)")
	flag.Parse()
}

func GetIPAndPort(ipPort string) (ip string, port uint16, err error) {
	nodeArr := strings.Split(ipPort, ":")
	if len(nodeArr) != 2 {
		err = fmt.Errorf("node ip and port format invalid. It should be ip:port\n")
		return
	}
	ip = nodeArr[0]
	portInt, err := strconv.Atoi(nodeArr[1])
	if err != nil {
		err = fmt.Errorf("could not convert port to int: %s\n", err)
		return
	}
	port = uint16(portInt)
	return
}

func main() {


	pk, sk, err := zmq4.NewCurveKeypair()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not generate curve key pair: %s\n", err)
		return
	}

	nodeIP, nodePort, err := GetIPAndPort(node)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	conf := config.Config{
		PublicKey:  pk,
		PrivateKey: sk,
		IP:         nodeIP,
		Port:       nodePort,
	}

		serverIP, serverPort, err := GetIPAndPort(server)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not parse ip and port of servers: %s\n", err)
			return
		}
		conf.Server = &config.ServerConfig{
			PublicKey: pk,
			IP:        serverIP,
			Port:      serverPort,
		}

	// write config
	v := viper.New()
	v.Set("config", conf)
	_, err = os.Stat(out)
	if !os.IsNotExist(err) {
		_, _ = fmt.Fprintf(os.Stderr, "error writing config file: file already exists\n")
		return
	}
	err = v.WriteConfigAs(out)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "cannot write config file: %s\n", err)
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, "config file written successfully in %s\n", out)
	_, _ = fmt.Fprintf(os.Stderr, "PUBLIC KEY: %s\n", pk)
}
