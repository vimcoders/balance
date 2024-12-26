package handler

import (
	"balance/pb"
	"context"
	"net"

	"github.com/vimcoders/go-driver/log"
	"github.com/vimcoders/go-driver/tcpx"
	//etcd "go.etcd.io/etcd/client/v3"
)

type Option struct {
	Balance struct {
		Address string `yaml:"address"`
	} `yaml:"balance"`
}

type Handler struct {
	Option
}

// MakeHandler creates a Handler instance
func MakeHandler(ctx context.Context) *Handler {
	h := &Handler{}
	return h
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, c net.Conn) {
	newSession := &Session{}
	cli := tcpx.NewClient(c, tcpx.Option{ServiceDesc: pb.Parkour_ServiceDesc})
	if err := cli.Register(ctx, newSession); err != nil {
		log.Error(err.Error())
	}
}

// Close stops handler
func (x *Handler) Close() error {
	return nil
}
