package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/BackAged/go-hexagonal-architecture/configuration"
	"github.com/BackAged/go-hexagonal-architecture/domain/task"
	"github.com/BackAged/go-hexagonal-architecture/infrastructure/database"
	"github.com/BackAged/go-hexagonal-architecture/infrastructure/repository"
	"github.com/go-chi/chi"
)

// Serve serves rest api
func Serve(cfgPath string) error {
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
	tskHndlr := NewHandler(tskSvc)

	r := chi.NewRouter()
	r.Mount("/api/v1/task", TaskRouter(tskHndlr))

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

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.GracefulTimeout)*time.Second)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
	return nil
}
