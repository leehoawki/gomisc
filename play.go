package main

import (
	"gomisc/log"
	"gomisc/tcp"
)

func main() {
	log.Setup("gomisc", "10.141.48.10:4560")
	tcp.ListenAndServeWithSignal(":8000", &tcp.ProxyHandler{})
}
