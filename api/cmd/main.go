package main

import (
	"api/internal/server"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)


func main() {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.TimeOnly})
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:9091", //env.MustHaveEnv(logger, "REDIS_URL"),
		Password: "root", //env.MustHaveEnv(logger, "REDIS_PASSWORD"),
		DB:       0,
	})

	defer rdb.Close()
	srv := http.Server{
		Addr: ":9999",
		// TLSConfig: &tls.Config{
		// 	GetCertificate: m.GetCertificate,
		// },
		Handler: server.NewRouter(logger, nil, rdb),
	}
	logger.Info().Msg("Server started")
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			logger.Error().Err(err).Msg("")
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	err := srv.Shutdown(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("")
	}
	cancel()
	logger.Info().Msg("Server gracefully shutted down")
}
