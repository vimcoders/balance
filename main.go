// 负载均衡服务 请求转发到服务去处理
// 同时我们也会主动推送很多数据，它的推送速度就像彗星一样的快
package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"balance/handler"

	"github.com/vimcoders/go-driver/log"
)

func main() {
	log.Info("NumCPU: ", runtime.NumCPU())
	quit := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	x := handler.MakeHandler(ctx)
	go x.ListenAndServe(ctx)
	log.Info("balance running")
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	s := <-quit
	switch s {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
		log.Info("shutdown ->", s.String())
		cancel()
		x.Close()
	default:
		log.Info("os.Signal ->", s.String())
	}
}
