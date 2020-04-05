package main

import (
	"context"
	"fmt"
	"github.com/BackAged/go-hexagonal-architecture/configuration"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start a http server",
	RunE:  serve,
}

func init() {
	serveCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "config.yaml", "config file path")
}

func serve(cmd *cobra.Command, args []string) error {
	cfg, err := configuration.Load(cfgPath)
	if err != nil {
		return err
	}

	// mgo, err := database.NewClient(cfg.Mongo.URI, cfg.Mongo.Database)
	// if err != nil {
	// 	return err
	// }

	// usrRepo := repo.NewMgoUser("users", mgo)
	// tagRepo := repo.NewMgoTag("tags", mgo)

	// svc := service.NewService(usrRepo, tagRepo)
	// r := svc.Route()
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "This is the people handler.")
		}),
	}

	go func() {
		log.Println("Staring server with address ", addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println("Failed to listen and serve", err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.GracefulTimeout)*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
	return nil
}
