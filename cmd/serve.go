package cmd

import (
	"github.com/niclabs/dhsm-signer/signer"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	runCmd.Flags().StringP("file", "f", "", "Full path to zone file to be verified")
	_ = runCmd.MarkFlagRequired("file")
}

var runCmd = &cobra.Command{
	Use:   "client",
	Short: "Runs the node",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
