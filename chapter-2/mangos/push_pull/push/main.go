package main

import (
	"go.nanomsg.org/mangos/v3/protocol/push"
	_ "go.nanomsg.org/mangos/v3/transport/tcp"
	"log"
	"time"
)

func main() {
	url := "tcp://127.0.0.1:40899"
	sock, err := push.NewSocket()
	if err != nil {
		log.Fatal("can't get new rep socket: %s", err.Error())
	}
	if err = sock.Dial(url); err != nil {
		log.Fatal("can't listen on rep socket: %s", err.Error())
	}

	err = sock.Send([]byte("Message 1"))
	if err != nil {
		log.Fatal("can't send message to rep socket: %s", err.Error())
	}
	time.Sleep(time.Second * 5)
	sock.Close()
}
