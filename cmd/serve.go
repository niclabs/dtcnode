package cmd

import (
	"github.com/niclabs/dtcnode/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs the node",
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Serve()
	},
}
