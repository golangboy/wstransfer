package main

import (
	"flag"
	"net"
	"net/url"

	"github.com/gorilla/websocket"
)

var serverAddr *string
var listenAddr *string

func handleConn(conn net.Conn) {
	defer conn.Close()
	u := url.URL{
		Scheme: "ws",
		Host:   *serverAddr,
		Path:   "/chat",
	}
	targetConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	defer targetConn.Close()
	go func() {
		for {
			bytes := make([]byte, 1024)
			n, err := conn.Read(bytes)
			if err != nil {
				return
			}
			targetConn.WriteMessage(websocket.BinaryMessage, bytes[:n])
		}
	}()
	for {
		_, message, err := targetConn.ReadMessage()
		if err != nil {
			return
		}
		conn.Write(message)
	}

}
func main() {
	serverAddr = flag.String("addr", "localhost:8080", "server address")
	listenAddr = flag.String("listen", "localhost:8081", "listen address")
	flag.Parse()
	listener, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}
