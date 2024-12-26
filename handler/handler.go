package handler

import (
	"balance/pb"
	"context"
	"net"
	"time"

	//etcd "go.etcd.io/etcd/client/v3"

	"github.com/vimcoders/go-driver/grpcx"
)

type Option struct {
	Balance struct {
		Address string `yaml:"address"`
	} `yaml:"balance"`
}

type Handler struct {
	Option
	grpcx.Client
}

// MakeHandler creates a Handler instance
func MakeHandler(ctx context.Context) *Handler {
	client, err := grpcx.Dial("udp", "127.0.0.1:19090", grpcx.Option{ServiceDesc: pb.Parkour_ServiceDesc})
	if err != nil {
		panic(err)
	}
	h := &Handler{Client: client}
	client.Register(ctx, h)
	return h
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, c net.Conn) {
	newSession := &Session{
		buffsize: 1024,
		timeout:  time.Second * 30,
		Conn:     c,
		Methods:  pb.Parkour_ServiceDesc.Methods,
		Client:   x.Client,
	}
	go newSession.serve(ctx)
	// cli := tcpx.NewClient(c, tcpx.Option{ServiceDesc: pb.Parkour_ServiceDesc})
	// if err := cli.Register(ctx, newSession); err != nil {
	// 	log.Error(err.Error())
	// }
}

// Close stops handler
func (x *Handler) Close() error {
	return nil
}
