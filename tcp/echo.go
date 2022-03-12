package tcp

/**
 * A echo server to test whether the server is functioning normally
 */

import (
	"bufio"
	"context"
	"gomisc/log"
	"io"
	"net"
	"sync"
)

// EchoHandler echos received line to client, using for test
type EchoHandler struct {
	activeConn sync.Map
}

// Handle echos received line to client
func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	h.activeConn.Store(conn, struct{}{})

	reader := bufio.NewReader(conn)
	for {
		// may occurs: client EOF, client timeout, server early close
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Info("connection close")
				h.activeConn.Delete(conn)
			} else {
				log.Error(err)
			}
			return
		}
		b := []byte(msg)
		_, _ = conn.Write(b)
	}
}

// Close stops echo handler
func (h *EchoHandler) Close() error {
	log.Info("handler shutting down...")
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(net.Conn)
		client.Close()
		return true
	})
	return nil
}
