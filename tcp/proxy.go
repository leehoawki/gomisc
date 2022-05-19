package tcp

import (
	"bufio"
	"context"
	"gomisc/log"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type ProxyHandler struct {
	activeConn sync.Map
}

type halfClosable interface {
	net.Conn
	CloseWrite() error
	CloseRead() error
}

func copyAndClose(dst, src halfClosable) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Error(err)
	}

	dst.CloseWrite()
	src.CloseRead()
}

// Handle echos received line to client
func (h *ProxyHandler) Handle(ctx context.Context, conn net.Conn) {
	h.activeConn.Store(conn, struct{}{})

	reader := bufio.NewReader(conn)
	var line, _, err = reader.ReadLine()
	segments := strings.Split(string(line), " ")
	if segments[0] != "CONNECT" {
		log.Info("not https request")
	}
	for {
		line, _, err = reader.ReadLine()
		data := (string)(line)
		println(data)
		if data == "" {
			break
		}
	}

	targetSiteCon, err := net.Dial("tcp", segments[1])
	targetSiteCon.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Error(err)
	}
	conn.Write([]byte("HTTP/1.0 200 Connection established\r\n\r\n"))
	targetTCP, targetOK := targetSiteCon.(halfClosable)
	proxyClientTCP, clientOK := conn.(halfClosable)

	if targetOK && clientOK {
		go copyAndClose(targetTCP, proxyClientTCP)
		go copyAndClose(proxyClientTCP, targetTCP)
	}
}

// Close stops echo handler
func (h *ProxyHandler) Close() error {
	log.Info("handler shutting down...")
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(net.Conn)
		client.Close()
		return true
	})
	return nil
}
