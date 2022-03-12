package tcp

/**
 * A tcp server
 */

import (
	"context"
	"fmt"
	"gomisc/log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type HandleFunc func(ctx context.Context, conn net.Conn)

// Handler represents application server over tcp
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}

// ListenAndServeWithSignal binds port and handle requests, blocking until receive stop signal
func ListenAndServeWithSignal(address string, handler Handler) error {
	closeChan := make(chan struct{})
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	//cfg.Address = listener.Addr().String()
	log.Info(fmt.Sprintf("bind: %s, start listening...", address))
	ListenAndServe(listener, handler, closeChan)
	return nil
}

// ListenAndServe binds port and handle requests, blocking until close
func ListenAndServe(listener net.Listener, handler Handler, closeChan <-chan struct{}) {
	// listen signal
	go func() {
		<-closeChan
		log.Info("shutting down...")
		_ = listener.Close() // listener.Accept() will return err immediately
		_ = handler.Close()  // close connections
	}()

	// listen port
	defer func() {
		// close during unexpected error
		_ = listener.Close()
		_ = handler.Close()
	}()
	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		// handle
		log.Info("accept link")
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
