package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/websocket"
)

func main() {
	host := flag.String("host", "0.0.0.0", "Host to register to")
	port := flag.String("port", "5000", "Master port number")
	id := flag.String("id", "", "Unique slave id (Required)")
	flag.Parse()

	// fmt.Println(*host, *port, *id)

	if len(*id) == 0 || len(*id) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	origin := "http://localhost/"
	url := fmt.Sprintf("ws://%s:%s/ws", *host, *port)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}
