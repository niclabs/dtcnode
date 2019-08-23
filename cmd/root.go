package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {
	rootCmd.AddCommand(genCurveCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(genConfigCmd)
	Log = log.New(os.Stderr, "", 0)
}

var Log *log.Logger

var rootCmd = &cobra.Command{
	Use:   "dtcnode",
	Short: "Executes a DTC Node",
	Long: `Executes a DTC Node.
	
	For more information, visit "https://github.com/niclabs/dtcnode".`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		Log.Printf("Error: %s", err)
		os.Exit(1)
	}
}
