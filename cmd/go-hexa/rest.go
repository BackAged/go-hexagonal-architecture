package main

import (
	"github.com/BackAged/go-hexagonal-architecture/application/rest"
	"github.com/spf13/cobra"
)

var serveRestCmd = &cobra.Command{
	Use:   "serve-rest",
	Short: "start a http server",
	RunE:  serve,
}

func init() {
	serveRestCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "config.yaml", "config file path")
}

func serve(cmd *cobra.Command, args []string) error {
	err := rest.Serve(cfgPath)
	return err
}
