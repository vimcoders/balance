package handler

import (
	"context"
	"net"

	"balance/pb"

	"github.com/vimcoders/go-driver/udpx"

	"github.com/vimcoders/go-driver/log"
	//etcd "go.etcd.io/etcd/client/v3"
)

type Handler struct {
	pb.UnimplementedParkourServer
	Option
}

// MakeHandler creates a Handler instance
func MakeHandler(ctx context.Context) *Handler {
	h := &Handler{}
	if err := h.Parse(); err != nil {
		panic(err)
	}
	if err := h.Connect(ctx); err != nil {
		panic(err)
	}
	return h
}

func (x *Handler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: req.Message}, nil
}

// Handle receives and executes redis commands
func (x *Handler) Handle(ctx context.Context, c net.Conn) {
	newSession := &Session{}
	cli := udpx.NewClient(c.(*net.UDPConn), udpx.Option{ServiceDesc: pb.Parkour_ServiceDesc})
	//cli := tcpx.NewClient(c, tcpx.Option{ServiceDesc: pb.Parkour_ServiceDesc})
	if err := cli.Register(ctx, newSession); err != nil {
		log.Error(err.Error())
	}
}

// Close stops handler
func (x *Handler) Close() error {
	return nil
}
