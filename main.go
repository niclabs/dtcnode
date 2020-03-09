package main

import (
	"fmt"
	"github.com/niclabs/dtcnode/v3/server"
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {
	Log = log.New(os.Stderr, "", 0)
}

var Log *log.Logger

func main() {
	viper.SetConfigName("dtcnode-config")
	viper.AddConfigPath("/etc/dtcnode/")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config file not found! %v", err))
	}
	if err := server.Serve(); err != nil {
		Log.Printf("Error: %s", err)
		os.Exit(1)
	}
}
