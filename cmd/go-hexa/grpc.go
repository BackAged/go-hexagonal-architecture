package main

import (
	"log"

	"github.com/spf13/cobra"
)

var serveGrpcCmd = &cobra.Command{
	Use:   "serve-grpc",
	Short: "start a grpc server",
	RunE:  serveGRPC,
}

func init() {
	serveGrpcCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "config.yaml", "config file path")
}

func serveGRPC(cmd *cobra.Command, args []string) error {
	log.Fatal("not implemented yet")
	return nil
}
