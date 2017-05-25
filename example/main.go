package main

import (
	"log"

	"time"

	"flag"

	"github.com/mushroomsir/engine.io"
)

func main() {
	for i := 0; i < 50000; i++ {
		go newClient()
		time.Sleep(1 * time.Millisecond)
		if i%1000 == 0 {
			log.Println("client count:", i)
		}
	}
	time.Sleep(time.Hour)
}

var ip = flag.String("ip", "127.0.0.1:0", "help message for flagname")

func newClient() {
	flag.Parse()
	conn, err := engineio.NewClient("ws://127.0.0.1:4000/engine.io/")
	if err != nil {
		log.Printf("newClient: %v", err)
		return
	}
	defer func() {
		conn.Close()
		log.Println("The client is closed")
	}()
	for {
		event := <-conn.Event
		switch event.Type {
		case "message":
			log.Println(string(event.Data))
		case "error":
			log.Println("Error:", string(event.Data))
		case "open":
			log.Println("open:", conn.GetSID())
		case "ping":
		default:
			//log.Println(event)
		}
	}
}
