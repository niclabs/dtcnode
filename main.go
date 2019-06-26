package main

import (
	"fmt"
	"github.com/niclabs/dtcnode/config"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"
	"os"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/dtcnode/")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config file not found! %v", err))
	}
	zmq4.AuthSetVerbose(true)
}

func main() {
	var conf config.Config
	err := viper.UnmarshalKey("config", &conf)
	if err != nil {
		panic(err)
	}
	if conf.Host == "" || conf.PublicKey == "" || conf.PrivateKey == "" || conf.Server == nil || conf.Port == 0 {
		panic(fmt.Errorf("missing fields in conf file"))
	}

	n, err := InitNode(&conf)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error initializing client: %s", err)
		return
	}
	n.Listen()
}
