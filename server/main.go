package main

import (
	"flag"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var targetAddr *string
var listenAddr *string

func main() {
	targetAddr = flag.String("addr", "localhost:8089", "server address")
	listenAddr = flag.String("listen", "localhost:80", "listen address")
	flag.Parse()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/chat", func(context *gin.Context) {
		upgrade := websocket.Upgrader{}
		conn, err := upgrade.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			return
		}
		targetConn, err := net.Dial("tcp", *targetAddr)
		if err != nil {
			return
		}
		go func() {
			defer conn.Close()
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					return
				}
				targetConn.Write(message)
			}
		}()
		go func() {
			defer targetConn.Close()
			for {
				bytes := make([]byte, 1024)
				n, err := targetConn.Read(bytes)
				if err != nil {
					return
				}
				conn.WriteMessage(websocket.BinaryMessage, bytes[:n])
			}
		}()

	})
	r.Run(*listenAddr)
}
