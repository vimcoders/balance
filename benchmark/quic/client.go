package quic

import (
	"balance/pb"
	"context"
	"crypto/tls"
	"time"

	"github.com/vimcoders/go-driver/quicx"
	"github.com/vimcoders/go-driver/tcpx"
)

type Client struct {
	tcpx.Client
}

func (x *Client) BenchmarkPing(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			x.Go(ctx, &pb.PingRequest{})
		}
	}
}

func Dail(address string) *Client {
	conn, err := quicx.Dial(address, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
		MaxVersion:         tls.VersionTLS13,
	}, &quicx.Config{
		MaxIdleTimeout: time.Minute,
	})
	if err != nil {
		panic(err)
	}
	return &Client{Client: tcpx.NewClient(conn, tcpx.Option{
		ServiceDesc: pb.Parkour_ServiceDesc,
	})}
}
