package genconfig

import (
	"fmt"
	"github.com/niclabs/dtcnode/v2/config"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
)

// GetHostAndPort splits a host and port string. Returns an error if something goes wrong.
func GetHostAndPort(ipPort string) (ip string, port uint16, err error) {
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

func GenerateConfig(node, client, pk, out string) error {
	_, err := os.Stat(out)
	if !os.IsNotExist(err) {
		_, _ = fmt.Fprintf(os.Stderr, "File already exists: Skipping config creation.\n")
		return nil
	}
	nodePK, nodeSK, err := zmq4.NewCurveKeypair()
	if err != nil {
		return fmt.Errorf("could not generate curve key pair: %s\n", err)
	}

	nodeHost, nodePort, err := GetHostAndPort(node)
	if err != nil {
		return fmt.Errorf("could not get host and port for the node: %s", err)
	}

	conf := config.Config{
		PublicKey:  nodePK,
		PrivateKey: nodeSK,
		Host:       nodeHost,
		Port:       nodePort,
	}

	if err != nil {
		return fmt.Errorf("could not parse ip and port of servers: %s\n", err)
	}

	conf.Client = &config.ClientConfig{
		PublicKey: pk,
		Host:      client,
	}

	v := viper.New()
	v.Set("config", conf)
	err = v.WriteConfigAs(out)
	if err != nil {
		return fmt.Errorf("cannot write config file: %s\n", err)
	}

	_, _ = fmt.Fprintf(os.Stderr, "config file written successfully in %s\n", out)
	_, _ = fmt.Fprintf(os.Stderr, "Public Key: %s\n", pk)
	return nil
}
