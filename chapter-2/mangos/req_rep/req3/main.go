package main

import (
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/req"
	"log"

	_ "go.nanomsg.org/mangos/v3/transport/tcp"
)

func sendMessage(sk mangos.Socket, msg string, ch chan bool) {
	ctx, err := sk.OpenContext()
	defer sk.Close()
	if err != nil {
		ch <- false
		return
	}
	err = ctx.Send([]byte(msg))
	if err != nil {
		ch <- false
		log.Fatal("can't send message to rep socket: %s", err.Error())
	}
	rsp, err := ctx.Recv()
	if err != nil {
		ch <- false
		log.Fatal("can't receive message from rep socket: %s", err.Error())
	}
	log.Println(string(rsp), err)
	ch <- true
}

func main() {
	url := "tcp://127.0.0.1:40899"
	sock, err := req.NewSocket()
	if err != nil {
		log.Fatal("can't get new rep socket: %s", err.Error())
	}
	if err = sock.Dial(url); err != nil {
		log.Fatal("can't listen on rep socket: %s", err.Error())
	}

	ch := make(chan bool, 3)

	go sendMessage(sock, "Message 1", ch)
	go sendMessage(sock, "Message 2", ch)
	go sendMessage(sock, "Message 3", ch)

	c := 0
	for rsp := range ch {
		log.Println(rsp)
		c += 1
		if c >= 3 {
			break
		}
	}
}
