package main

import (
	"go.nanomsg.org/mangos/v3/protocol/req"
	"log"

	_ "go.nanomsg.org/mangos/v3/transport/tcp"
)

func main() {
	url := "tcp://127.0.0.1:40899"
	sock, err := req.NewSocket()
	if err != nil {
		log.Fatal("can't get new rep socket: %s", err.Error())
	}
	if err = sock.Dial(url); err != nil {
		log.Fatal("can't listen on rep socket: %s", err.Error())
	}

	ctx1, err := sock.OpenContext()
	if err != nil {
		log.Fatal("can't create context: %s", err.Error())
	}

	ctx2, err := sock.OpenContext()
	if err != nil {
		log.Fatal("can't create context: %s", err.Error())
	}

	err = ctx1.Send([]byte("Message 1"))
	if err != nil {
		log.Fatal("can't send message to rep socket: %s", err.Error())
	}

	err = ctx2.Send([]byte("Message 2"))
	if err != nil {
		log.Fatal("can't send message to rep socket: %s", err.Error())
	}
	msg, err := ctx1.Recv()
	log.Println(string(msg), err)
	msg, err = ctx2.Recv()
	log.Println(string(msg), err)
}
