package main

import (
	"fmt"
	config2 "github.com/niclabs/dtcnode/config"
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
}


func main() {
	var config config2.Config
	err := viper.UnmarshalKey("dtcnode", &config)
	if err != nil {
		panic(err)
	}

	n, err := InitClient(&config)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error initializing client: %s", err)
		return
	}
	n.Listen()
}
