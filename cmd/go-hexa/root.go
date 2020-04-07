package main

import (
	"github.com/spf13/cobra"
)

// Version defines version
const (
	Version = "unversioned"
)

// rootCmd is the root of all sub commands in the binary
// it doesn't have a Run method as it executes other sub commands
var rootCmd = &cobra.Command{
	Use:     "user task",
	Short:   "task manages user task",
	Version: Version,
}

// for now only requirement is to run server
func main() {
	rootCmd.Execute()
}

var cfgPath string

func init() {
	// register sub-commands here
	rootCmd.AddCommand(serveRestCmd)
	rootCmd.AddCommand(serveGrpcCmd)
}
