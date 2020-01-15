package server

import (
	"fmt"
	"github.com/niclabs/dtcnode/v3/config"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"
)

func Serve() error {
	var conf config.Config
	err := viper.UnmarshalKey("config", &conf)
	if err != nil {
		return err
	}
	if conf.Host == "" {
		conf.Host = "0.0.0.0"
	}
	if conf.PublicKey == "" || conf.PrivateKey == "" || conf.Client == nil || conf.Port == 0 {
		return fmt.Errorf("missing fields in conf file")
	}

	err = zmq4.AuthStart()
	if err != nil {
		return fmt.Errorf("error starting auth: %s", err)
	}
	defer zmq4.AuthStop()

	n, err := InitNode(&conf)
	if err != nil {
		return fmt.Errorf("error initializing node: %s", err)
	}
	n.Listen()
	return nil
}