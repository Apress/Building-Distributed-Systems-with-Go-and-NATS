package main

import (
	"fmt"
	"go.nanomsg.org/mangos/v3/protocol/pub"
	_ "go.nanomsg.org/mangos/v3/transport/tcp"
	"log"
	"time"
)

func main() {
	tt := []string{"Blue - ", "Red - ", "Yellow - "}
	url := "tcp://127.0.0.1:40899"
	sock, err := pub.NewSocket()
	if err != nil {
		log.Fatal("can't get new push socket: %s", err.Error())
	}
	if err = sock.Listen(url); err != nil {
		log.Fatal("can't listen on push socket: %s", err.Error())
	}

	i := 0
	for {
		msg := fmt.Sprintf("%s Message %d", tt[i%3], i)
		i += 1
		err = sock.Send([]byte(msg))
		if err != nil {
			log.Fatal("can't send message to pull socket: %s", err.Error())
		}
		log.Println("Sent message", msg)
		time.Sleep(time.Second * 5)
	}

}
