package main

import (
	"balance/benchmark/quic"
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/vimcoders/go-driver/log"
)

func main() {
	log.Info("NumCPU: ", runtime.NumCPU())
	x := quic.Dail("127.0.0.1:9090")
	go x.ListenAndServe(context.TODO())
	go x.BenchmarkLogin(context.Background())
	quit := make(chan os.Signal, 1)
	log.Info("balance running")
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	s := <-quit
	switch s {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
		log.Info("shutdown ->", s.String())
	default:
		log.Info("os.Signal ->", s.String())
	}
}
