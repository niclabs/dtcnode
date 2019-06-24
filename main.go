package main

import (
	"fmt"
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
	var config Config
	err := viper.UnmarshalKey("config", &config)
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
