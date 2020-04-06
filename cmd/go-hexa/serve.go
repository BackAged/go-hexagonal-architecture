package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/BackAged/go-hexagonal-architecture/application/rest"
	"github.com/BackAged/go-hexagonal-architecture/configuration"
	"github.com/BackAged/go-hexagonal-architecture/domain/task"
	"github.com/BackAged/go-hexagonal-architecture/infrastructure/database"
	"github.com/BackAged/go-hexagonal-architecture/infrastructure/repository"
	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
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

	rds, err := database.NewInMemoryClient(cfg.Redis.Host, cfg.Redis.Password, &cfg.Redis.DB)
	if err != nil {
		return err
	}

	tskRepo := repository.NewTaskRepository(rds)
	tskSvc := task.NewService(tskRepo)
	tskHndlr := rest.NewHandler(tskSvc)

	r := chi.NewRouter()
	r.Mount("/task", rest.TaskRouter(tskHndlr))
	// r := svc.Route()
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler:      r,
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
