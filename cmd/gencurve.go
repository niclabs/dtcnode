package cmd

import (
	"fmt"
	"github.com/pebbe/zmq4"
	"github.com/spf13/cobra"
)

var genCurveCmd = &cobra.Command{
	Use:   "generate-curve",
	Short: "Returns a Public Key and a Private Key usable in the system",
	RunE: func(cmd *cobra.Command, args []string) error {
		pk, sk, err := zmq4.NewCurveKeypair()
		if err != nil {
			return fmt.Errorf("Cannot generate curve key pair")
		}
		fmt.Printf("Public Key:   %s\n", pk)
		fmt.Printf("Private Key:  %s\n", sk)
		return nil
	},
}
