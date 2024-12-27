package handler

import (
	"balance/pb"
	"context"
	"net"
	"time"

	//etcd "go.etcd.io/etcd/client/v3"

	"github.com/vimcoders/go-driver/grpcx"
	"google.golang.org/protobuf/proto"
)

type Option struct {
	Balance struct {
		Address string `yaml:"address"`
	} `yaml:"balance"`
}

type Handler struct {
	Methods []Method
	Option
	grpcx.Client
}

// MakeHandler creates a Handler instance
func MakeHandler(ctx context.Context) *Handler {
	var methods []Method
	for i := 0; i < len(pb.Parkour_ServiceDesc.Methods); i++ {
		var newMethod Method
		method := pb.Parkour_ServiceDesc.Methods[i]
		dec := func(in any) error {
			newMethod.Request = in.(proto.Message)
			return nil
		}
		resp, _ := method.Handler(&pb.UnimplementedParkourServer{}, context.Background(), dec, nil)
		newMethod.Id = uint16(len(methods) + 512)
		newMethod.ServiceName = pb.Parkour_ServiceDesc.ServiceName
		newMethod.MethodName = method.MethodName
		newMethod.Reply = resp.(proto.Message)
		methods = append(methods, newMethod)
	}
	for i := 0; i < len(pb.Chat_ServiceDesc.Methods); i++ {
		var newMethod Method
		method := pb.Chat_ServiceDesc.Methods[i]
		dec := func(in any) error {
			newMethod.Request = in.(proto.Message)
			return nil
		}
		resp, _ := method.Handler(&pb.UnimplementedChatServer{}, context.Background(), dec, nil)
		newMethod.Id = uint16(len(methods) + 512)
		newMethod.ServiceName = pb.Chat_ServiceDesc.ServiceName
		newMethod.MethodName = method.MethodName
		newMethod.Reply = resp.(proto.Message)
		methods = append(methods, newMethod)
	}
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
		Client:   x.Client,
	}
	for i := 0; i < len(x.Methods); i++ {
		newSession.Methods = append(newSession.Methods, x.Methods[i].Clone())
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
