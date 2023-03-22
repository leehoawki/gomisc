package tcp

import (
	"bufio"
	"context"
	"gomisc/log"
	"io/ioutil"
	"net"
	"strings"
	"sync"
)

type HttpHandler struct {
	activeConn sync.Map
}

func cat(filename string) []byte {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return notFound()
	}
	return content
}

func notFound() []byte {
	return []byte("HTTP/1.0 404 NOT FOUND\r\n" +
		"Server: gohttpd/0.1.0\r\n" +
		"Content-Type: text/html\r\n" +
		"\r\n" +
		"<HTML><TITLE>Not Found</TITLE>\r\n" +
		"<BODY><P>The server could not fulfill\r\n" +
		"your request because the resource specified\r\n" +
		"is unavailable or nonexistent.\r\n" +
		"</BODY></HTML>\r\n")
}

func badRequest() []byte {
	return []byte("HTTP/1.0 400 BAD REQUEST\r\n" +
		"Content-type: text/html\r\n" +
		"\r\n" +
		"<P>Your browser sent a bad request, " +
		"such as a POST without a Content-Length.\r\n")
}

func headers() []byte {
	return []byte("HTTP/1.0 200 OK\r\n" +
		"Server: gohttpd/0.1.0\r\n" +
		"Content-Type: text/html\r\n" +
		"\r\n")
}

func (h *HttpHandler) Handle(ctx context.Context, conn net.Conn) {
	h.activeConn.Store(conn, struct{}{})

	reader := bufio.NewReader(conn)
	var line, _, err = reader.ReadLine()
	if err != nil {
		conn.Write(badRequest())
		h.closeConnect(conn)
		return
	}

	segments := strings.Split(string(line), " ")
	method := segments[0]
	path := segments[1]
	log.Info("request accept, method=" + method + ", path=" + path)
	if "GET" != method {
		conn.Write(badRequest())
		h.closeConnect(conn)
		return
	}
	conn.Write(headers())
	conn.Write(cat("." + path))
	h.closeConnect(conn)
}

func (h *HttpHandler) closeConnect(conn net.Conn) {
	h.activeConn.Delete(conn)
	conn.Close()
}

func (h *HttpHandler) Close() error {
	log.Info("handler shutting down...")
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(net.Conn)
		client.Close()
		return true
	})
	return nil
}
