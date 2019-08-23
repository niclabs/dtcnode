package cmd

import (
	"github.com/niclabs/dtcnode/genconfig"
	"github.com/spf13/cobra"
)

var node string
var client string
var pk string
var out string

func init() {
	genConfigCmd.Flags().StringVarP(&node, "node", "n", "0.0.0.0:2030", "node ip and port")
	genConfigCmd.Flags().StringVarP(&client, "client", "c", "", "client ip")
	genConfigCmd.Flags().StringVarP(&pk, "key", "k", "", "base85 client public key")
	genConfigCmd.Flags().StringVarP(&out, "output", "o", "", "path where to output the config file (default: ./config.yaml)")
	_ = genConfigCmd.MarkFlagRequired("client")
	_ = genConfigCmd.MarkFlagRequired("key")
	_ = genConfigCmd.MarkFlagRequired("output")

}

var genConfigCmd = &cobra.Command{
	Use:   "generate-config",
	Short: "Generates a configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return genconfig.GenerateConfig(node, client, pk, out)
	},
}
