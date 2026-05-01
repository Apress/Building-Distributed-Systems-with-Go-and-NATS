package main

import (
	"fmt"
	"go.nanomsg.org/mangos/v3/protocol/push"
	_ "go.nanomsg.org/mangos/v3/transport/tcp"
	"log"
	"time"
)

func main() {
	url := "tcp://127.0.0.1:40899"
	sock, err := push.NewSocket()
	if err != nil {
		log.Fatal("can't get new push socket: %s", err.Error())
	}
	if err = sock.Listen(url); err != nil {
		log.Fatal("can't listen on push socket: %s", err.Error())
	}

	i := 0
	for {
		i += 1
		err = sock.Send([]byte(fmt.Sprintf("Message %d", i)))
		if err != nil {
			log.Fatal("can't send message to pull socket: %s", err.Error())
		}
		log.Println("Sent message")
		time.Sleep(time.Second * 5)
	}

}
