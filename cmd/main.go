package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/wonwooseo/panawa-api/router"
)

func main() {
	cfgF := flag.String("config", "", "path to config file")
	portF := flag.Int("port", 80, "port number to listen(default: 80)")
	flag.Parse()

	baseLogger := log.Logger
	logger := baseLogger.With().Str("caller", "main").Logger()

	viper.SetConfigFile(*cfgF)
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal().Err(err).Msg("failed to read in config")
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", *portF),
		Handler:      router.NewRouter(baseLogger),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	serverCtx, serverCancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info().Msg("got signal")

		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 10*time.Second)
		go func() {
			defer shutdownCancel()
			<-shutdownCtx.Done() // if serverCtx is not cancelled before shutdownCtx times out, forcefully exit
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				logger.Fatal().Msg("graceful shutdown timeout; exiting forcefully")
			}
		}()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Fatal().Err(err).Msg("graceful shutdown failed; exiting forcefully")
		}
		serverCancel()
	}()

	logger.Info().Str("addr", srv.Addr).Msg("start listening..")
	if err := srv.ListenAndServe(); err != nil {
		logger.Error().Err(err).Msg("failed to start listening")
	}

	<-serverCtx.Done() // wait for graceful shutdown
}
