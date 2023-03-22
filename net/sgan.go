package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var C1 uint32 = 0x1000000
var C2 uint32 = 0x10000
var C3 uint32 = 0x100

func fromIP(ip uint32) string {
	var i = ip
	s1 := ip / C1
	i = i - s1*C1
	s2 := i / C2
	i = i - s2*C2
	s3 := i / C3
	i = i - s3*C3
	s4 := i

	return strconv.Itoa(int(s1)) + "." +
		strconv.Itoa(int(s2)) + "." +
		strconv.Itoa(int(s3)) + "." +
		strconv.Itoa(int(s4))
}

func toIP(ip string) uint32 {
	segments := strings.Split(ip, ".")
	var x uint32 = 0
	for _, segment := range segments {
		x = x * C3
		val, err := strconv.Atoi(segment)
		if err != nil {
			panic(err)
		}
		x += uint32(val)
	}
	return x
}

func main() {
	parser := argparse.NewParser("sgan", "go version port scanner")
	s := parser.String("s", "start", &argparse.Options{Required: false, Help: "start ip"})
	e := parser.String("e", "end", &argparse.Options{Required: false, Help: "end ip"})
	t := parser.String("t", "target", &argparse.Options{Required: false, Help: "target ip"})
	p := parser.Int("p", "port", &argparse.Options{Required: false, Help: "target port when scanning ip range "})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(-1)
	}
	start := *s
	end := *e
	target := *t
	port := *p
	if target != "" {
		println("scanning target, ip=" + target)
		var wg sync.WaitGroup
		channel := make(chan int, 0xffff)
		for i := 1; i < runtime.NumCPU()*256; i++ {
			go func() {
				for {
					j := <-channel
					if doScan(target, j) {
						println("port open:" + strconv.Itoa(j))
					}
					wg.Done()
				}
			}()
		}
		for i := 1; i <= 0xffff; i++ {
			channel <- i
			wg.Add(1)
		}
		wg.Wait()
	} else if start != "" && end != "" {
		println("scanning port, target=" + strconv.Itoa(port) + ", from=" + start + ", to=" + end)
		s := toIP(start)
		e := toIP(end)
		if s > e {
			fmt.Print(parser.Usage(err))
			os.Exit(-1)
		}
		var wg sync.WaitGroup
		channel := make(chan uint32, e-s+1)
		for i := 1; i < runtime.NumCPU()*256; i++ {
			go func() {
				for {
					j := <-channel
					address := fromIP(j)
					if doScan(address, port) {
						println("ip:" + address)
					}
					wg.Done()
				}
			}()
		}
		for i := s; i <= e; i++ {
			channel <- i
			wg.Add(1)
		}
		wg.Wait()
	} else {
		fmt.Print(parser.Usage(err))
		os.Exit(-1)
	}
}

func doScan(ip string, port int) bool {
	_, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err == nil {
		return true
	}
	return false
}
