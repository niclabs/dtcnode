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
	flag.StringVar(&server, "s", "", "server ip and port, separated by commas if more than one server")
	flag.StringVar(&pk, "k", "", "base85 server public key, separated by commas if more than one server")
	flag.StringVar(&out, "o", "./config.yaml", "path where to output the config file (default: ./config.yaml)")
	flag.Parse()
}

func GetIPAndPort(ipPort string) (ip string, port uint16, err error) {
	nodeArr := strings.Split(ipPort, ":")
	if len(nodeArr) != 2 {
		err = fmt.Errorf("node ip and port format invalid. It should be ip:port\n")
		return
	}
	portInt, err := strconv.Atoi(nodeArr[1])
	if err != nil {
		err = fmt.Errorf("could not convert port to int: %s\n", err)
		return
	}
	port = uint16(portInt)
	return
}

func main() {

	servers := strings.Split(server, ",")
	pks := strings.Split(pk, ",")

	if len(servers) != len(pks) {
		_, _ = fmt.Fprintf(os.Stderr, "number of servers, pks and sks provided is different\n")
		return
	}

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
		Servers:    make([]*config.ServerConfig, len(servers)),
	}

	for i := 0; i < len(servers); i++ {
		serverIP, serverPort, err := GetIPAndPort(servers[i])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not parse ip and port of servers: %s\n", err)
			return
		}
		conf.Servers[i] = &config.ServerConfig{
			PublicKey: pks[i],
			IP:        serverIP,
			Port:      serverPort,
		}
	}

	// write config
	v := viper.New()
	v.Set("config", conf)
	err = v.WriteConfigAs(out)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error writing config file: %s\n", err)
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, "config file written successfully in %s\n", out)
	_, _ = fmt.Fprintf(os.Stderr, "PUBLIC KEY: %s\n", pk)
}