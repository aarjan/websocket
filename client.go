package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var host = flag.String("host", "0.0.0.0", "Host to register to")
var port = flag.String("port", "5000", "Master port number")
var id = flag.String("id", "", "Unique slave id (Required)")

func main() {
	flag.Parse()

	if len(*id) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	log.SetFlags(0)

	// For user interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	addr := fmt.Sprintf("%s:%s", *host, *port)
	u := url.URL{Scheme: "ws", Host: addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// dummy channel to end the session
	done := make(chan int)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	// Ticker with a duration of 1 sec.
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		// Msg sent with a recieve from ticker channel
		case t := <-ticker.C:
			// Current time is sent as a text message
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "connection interrupted"))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
		}
	}
}
