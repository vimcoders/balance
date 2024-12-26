package handler

import (
	"balance/pb"
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/vimcoders/go-driver/grpcx"
	"github.com/vimcoders/go-driver/log"
	"google.golang.org/protobuf/proto"
)

type Method struct {
	MethodName string
	Request    proto.Message
	Reply      proto.Message
}

func (x Method) Clone() *Method {
	return &Method{
		MethodName: x.MethodName,
		Request:    x.Request.ProtoReflect().New().Interface(),
		Reply:      x.Reply.ProtoReflect().New().Interface(),
	}
}

type Session struct {
	net.Conn
	grpcx.Client
	buffsize int
	timeout  time.Duration
	Methods  []Method
	pb.UnimplementedParkourServer
}

func (x *Session) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	log.Debug("Ping")
	return &pb.PingResponse{Message: req.Message}, nil
}

func (x *Session) Close() error {
	return nil
}

func (x *Session) serve(ctx context.Context) (err error) {
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
		if err := x.Conn.SetReadDeadline(time.Now().Add(x.timeout)); err != nil {
			return err
		}
		iMessage, err := decode(buf)
		if err != nil {
			return err
		}
		if x == nil {
			continue
		}
		method := x.Methods[iMessage.method()]
		methodName, req, reply := method.MethodName, method.Request, method.Reply
		if err := proto.Unmarshal(iMessage.payload(), method.Request); err != nil {
			return err
		}
		if err := x.Invoke(ctx, methodName, req, reply); err != nil {
			return err
		}
		b, err := encode(iMessage.method(), method.Reply)
		if err != nil {
			return err
		}
		if _, err := x.Write(b); err != nil {
			return err
		}
	}
}

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

func (x Message) length() uint16 {
	return binary.BigEndian.Uint16(x)
}

func (x Message) method() uint16 {
	return binary.BigEndian.Uint16(x[2:])
}

func (x Message) payload() []byte {
	return x[4:x.length()]
}

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
