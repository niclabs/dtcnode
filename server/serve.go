package server

import (
	"flag"
	"fmt"
	"github.com/niclabs/dtcnode/config"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"
	"os"
)

var genConfig bool

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/dtcnode/")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config file not found! %v", err))
	}
	genConfig = *flag.Bool("g", false, "If set, the program generates a config using")
}

func main() {
	var conf config.Config
	err := viper.UnmarshalKey("config", &conf)
	if err != nil {
		panic(err)
	}
	if conf.Host == "" {
		conf.Host = "0.0.0.0"
	}
	if conf.PublicKey == "" || conf.PrivateKey == "" || conf.Client == nil || conf.Port == 0 {
		panic(fmt.Errorf("missing fields in conf file"))
	}
	err = zmq4.AuthStart()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error starting auth: %s", err)
		return
	}
	n, err := InitNode(&conf)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error initializing node: %s", err)
		return
	}
	defer zmq4.AuthStop()
	n.Listen()
}
