package main

import (
	"fmt"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./") // config should be in the same folder
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config file not found! %v", err))
	}
}


func main() {
	// Start node
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	for i, server := range config.Servers {
		serverID := fmt.Sprintf("server_%d", i)
		zmq4.AuthAllow(serverID, server.IP)
		zmq4.AuthCurveAdd(serverID, server.PublicKey)
	}

	n, err := NewNode(config.PublicKey, config.PrivateKey, config.IP, config.RouterPort, config.SubPort)
	if err != nil {
		panic(err)
	}

	for {
		n.Listen()
	}
}