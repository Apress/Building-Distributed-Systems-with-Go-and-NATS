package main

import (
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/pull"
	"go.nanomsg.org/mangos/v3/protocol/push"
	_ "go.nanomsg.org/mangos/v3/transport/tcp"
	"log"
)

func handleMessage(ch chan []byte, sk mangos.Socket, id int) {
	for {
		msg := <-ch
		log.Println("worker ", id, " Working on message ", string(msg))
		sk.Send([]byte("Sending back: " + string(msg)))
	}

}

func main() {

	ch := make(chan []byte)
	plhurl := "tcp://127.0.0.1:40898"
	plsock, err := pull.NewSocket()
	if err != nil {
		log.Fatal("can't get new rep socket: %s", err.Error())
	}
	if err = plsock.Listen(plhurl); err != nil {
		log.Fatal("can't listen on rep socket: %s", err.Error())
	}

	psurl := "tcp://127.0.0.1:40899"
	pssock, err := push.NewSocket()
	if err != nil {
		log.Fatal("can't get new rep socket: %s", err.Error())
	}
	if err = pssock.Listen(psurl); err != nil {
		log.Fatal("can't listen on rep socket: %s", err.Error())
	}

	go handleMessage(ch, pssock, 1)
	go handleMessage(ch, pssock, 2)
	for {
		msg, err := plsock.Recv()
		if err != nil {
			log.Fatal("cannot receive on pull socket: %s", err.Error())
		}
		log.Println(string(msg))
		ch <- msg
	}

}
