package reverse_proxy

import (
	"context"
	"fmt"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/keepalive"
	"math"
	"net"
	"strings"
	"sync"
)

const (
	defaultServerMaxReceiveMessageSize = 1024 * 1024 * 4
	defaultServerMaxSendMessageSize    = math.MaxInt32
)

func ServeTCP(ctx context.Context, c net.Conn) {
	//rawConn.SetDeadline(time.Now().Add(s.opts.connectionTimeout))
	//conn, authInfo, err := s.useTransportAuthenticator(rawConn)
	//return rawConn, nil, nil
	//st := s.newHTTP2Transport(conn, authInfo)
	config := &transport.ServerConfig{
		MaxStreams:            0,
		AuthInfo:              nil,
		InTapHandle:           nil,
		StatsHandler:          nil,
		KeepaliveParams:       keepalive.ServerParameters{},
		KeepalivePolicy:       keepalive.EnforcementPolicy{},
		InitialWindowSize:     0,
		InitialConnWindowSize: 0,
		WriteBufferSize:       defaultServerMaxSendMessageSize,
		ReadBufferSize:        defaultServerMaxSendMessageSize,
		ChannelzParentID:      0,
		MaxHeaderListSize:     nil,
		HeaderTableSize:       nil,
	}
	st, err := transport.NewServerTransport("http2", c, config)
	if err != nil {
		fmt.Println("NewServerTransport(%q) failed: %v", c.RemoteAddr(), err)
		//c.Close()
		return
	}
	//rawConn.SetDeadline(time.Time{})
	//s.serveStreams(st)
	defer st.Close()
	var wg sync.WaitGroup
	st.HandleStreams(func(stream *transport.Stream) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sm := stream.Method()
			if sm != "" && sm[0] == '/' {
				sm = sm[1:]
			}
			pos := strings.LastIndex(sm, "/")
			if pos == -1 {
				errDesc := fmt.Sprintf("malformed method name: %q", stream.Method())
				fmt.Println(errDesc)
				return
			}
			service := sm[:pos]
			method := sm[pos+1:]
			fmt.Println("method", method)
			fmt.Println("service", service)
			//s.handleStream(st, stream, s.traceInfo(st, stream))
		}()
	}, func(ctx context.Context, method string) context.Context {
		return ctx
		//tr := trace.New("grpc.Recv."+methodFamily(method), method)
		//return trace.NewContext(ctx, tr)
	})
	wg.Wait()
	//if !s.addConn(st) {
	//	return
	//}
	//s.serveStreams(st)
}