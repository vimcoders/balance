package quic

import (
	"balance/pb"
	"bufio"
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/vimcoders/go-driver/log"
	"github.com/vimcoders/go-driver/quicx"
	"google.golang.org/protobuf/proto"
)

type Message []byte

var pool sync.Pool = sync.Pool{
	New: func() any {
		return &Message{}
	},
}

func decode(b *bufio.Reader) (Message, error) {
	headerBytes, err := b.Peek(2)
	if err != nil {
		return nil, err
	}
	length := int(binary.BigEndian.Uint16(headerBytes))
	if length > b.Size() {
		return nil, fmt.Errorf("header %v too long", length)
	}
	iMessage, err := b.Peek(length)
	if err != nil {
		return nil, err
	}
	if _, err := b.Discard(len(iMessage)); err != nil {
		return nil, err
	}
	return iMessage, nil
}

func encode(method uint16, iMessage proto.Message) (Message, error) {
	b, err := proto.Marshal(iMessage)
	if err != nil {
		return nil, err
	}
	buf := pool.Get().(*Message)
	buf.WriteUint16(uint16(4 + len(b)))
	buf.WriteUint16(method)
	if _, err := buf.Write(b); err != nil {
		return nil, err
	}
	return *buf, nil
}

// func (x Message) length() uint16 {
// 	return binary.BigEndian.Uint16(x)
// }

// func (x Message) method() uint16 {
// 	return binary.BigEndian.Uint16(x[2:])
// }

// func (x Message) payload() []byte {
// 	return x[4:x.length()]
// }

func (x *Message) reset() {
	if cap(*x) <= 0 {
		return
	}
	*x = (*x)[:0]
	pool.Put(x)
}

func (x *Message) Write(p []byte) (int, error) {
	*x = append(*x, p...)
	return len(p), nil
}

func (x *Message) WriteUint32(v uint32) {
	*x = binary.BigEndian.AppendUint32(*x, v)
}

func (x *Message) WriteUint16(v uint16) {
	*x = binary.BigEndian.AppendUint16(*x, v)
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
func (x Message) WriteTo(w io.Writer) (n int64, err error) {
	if nBytes := len(x); nBytes > 0 {
		m, e := w.Write(x)
		if m > nBytes {
			panic("bytes.Buffer.WriteTo: invalid Write count")
		}
		if e != nil {
			return n, e
		}
		// all bytes should have been written, by definition of
		// Write method in io.Writer
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	// Buffer is now empty; reset.
	x.reset()
	return n, nil
}

type Method struct {
	Id         uint16
	MethodName string
	Request    proto.Message
	Reply      proto.Message
}

func (x Method) Clone() Method {
	return Method{
		Id:         x.Id,
		MethodName: x.MethodName,
		Request:    x.Request.ProtoReflect().New().Interface(),
		Reply:      x.Reply.ProtoReflect().New().Interface(),
	}
}

type Client struct {
	net.Conn
	buffsize int
	timeout  time.Duration
	Methods  []Method
}

func (x *Client) Go(ctx context.Context, metodName string, req proto.Message) (err error) {
	defer func() {
		if err != nil {
			log.Error(err)
			debug.PrintStack()
		}
	}()
	for i := 0; i < len(x.Methods); i++ {
		if ok := strings.EqualFold(metodName, x.Methods[i].MethodName); !ok {
			continue
		}
		if err := x.push(x.Methods[i].Id, req); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("%s not registed", metodName)
}

func (x *Client) push(method uint16, req proto.Message) (err error) {
	buf, err := encode(method, req)
	if err != nil {
		return err
	}
	if err := x.SetWriteDeadline(time.Now().Add(x.timeout)); err != nil {
		return err
	}
	if _, err := buf.WriteTo(x.Conn); err != nil {
		return err
	}
	return nil
}

func (x *Client) BenchmarkLogin(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			log.Info("BenchmarkLogin")
			x.Go(ctx, "Login", &pb.LoginRequest{Token: "token"})
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
	var methods []Method
	for i := 0; i < len(pb.Parkour_ServiceDesc.Methods); i++ {
		var newMethod Method
		method := pb.Parkour_ServiceDesc.Methods[i]
		resp, _ := method.Handler(&pb.UnimplementedParkourServer{}, context.Background(), func(in any) error {
			newMethod.Request = in.(proto.Message)
			return nil
		}, nil)
		newMethod.Id = uint16(len(methods) + 512)
		newMethod.MethodName = method.MethodName
		newMethod.Reply = resp.(proto.Message)
		methods = append(methods, newMethod)
	}
	for i := 0; i < len(pb.Chat_ServiceDesc.Methods); i++ {
		var newMethod Method
		method := pb.Chat_ServiceDesc.Methods[i]
		resp, _ := method.Handler(&pb.UnimplementedChatServer{}, context.Background(), func(in any) error {
			newMethod.Request = in.(proto.Message)
			return nil
		}, nil)
		newMethod.Id = uint16(len(methods) + 512)
		newMethod.MethodName = method.MethodName
		newMethod.Reply = resp.(proto.Message)
		methods = append(methods, newMethod)
	}
	return &Client{
		Conn:     conn,
		timeout:  time.Minute,
		buffsize: 1024,
		Methods:  methods,
	}
}

func (x *Client) ListenAndServe(ctx context.Context) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
		if err != nil {
			log.Error(err.Error())
			debug.PrintStack()
		}
		if err := x.Close(); err != nil {
			log.Error(err.Error())
		}
	}()
	buf := bufio.NewReaderSize(x.Conn, x.buffsize)
	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		default:
		}
		log.Debug(x.timeout)
		if err := x.Conn.SetReadDeadline(time.Now().Add(x.timeout)); err != nil {
			return err
		}
		_, err := decode(buf)
		if err != nil {
			return err
		}
		// log.Debug(iMessage)
		// if x == nil {
		// 	continue
		// }
		// method, payload := iMessage.method(), iMessage.payload()
		// dec := func(in any) error {
		// 	if err := proto.Unmarshal(payload, in.(proto.Message)); err != nil {
		// 		return err
		// 	}
		// 	return nil
		// }
		// reply, err := x.Methods[method].Handler(x, ctx, dec, nil)
		// if err != nil {
		// 	return err
		// }
		// // 从网关转发到远程接口调用
		// // if err := x.Invoke(ctx, x.Methods[method].MethodName, req, reply); err != nil {
		// // 	return err
		// // }
		// b, err := encode(method, reply.(proto.Message))
		// if err != nil {
		// 	return err
		// }
		// if _, err := x.Write(b); err != nil {
		// 	return err
		// }
	}
}
