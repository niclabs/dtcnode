package cmd

import (
	"fmt"
	"github.com/niclabs/dtcnode/v2/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs the node",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/dtcnode/")
		viper.AddConfigPath("./")
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("config file not found! %v", err))
		}
		return server.Serve()
	},
}
