package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/wonwooseo/panawa-api/router"
)

var (
	Version   = ""
	BuildTime = ""
)

func main() {
	logger := log.Logger.With().Str("caller", "main").Logger()

	srv := &http.Server{
		Addr:         ":80",
		Handler:      router.NewRouter(Version, BuildTime),
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
