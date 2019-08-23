package main

import (
	"fmt"
	"github.com/niclabs/dtcnode/cmd"
	"github.com/spf13/viper"
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
	cmd.Execute()
}
