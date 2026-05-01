package main

import (
	"fmt"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/pub"
	"go.nanomsg.org/mangos/v3/protocol/sub"
	_ "go.nanomsg.org/mangos/v3/transport/inproc"
	"log"
	"time"
)

func publisher() {
	url := "inproc://publisher"
	sock, err := pub.NewSocket()
	if err != nil {
		log.Fatal("can't get new push socket: %s", err.Error())
	}
	if err = sock.Listen(url); err != nil {
		log.Fatal("can't listen on push socket: %s", err.Error())
	}

	time.Sleep(time.Second * 5)
	i := 0
	for {
		i += 1
		err = sock.Send([]byte(fmt.Sprintf("Message %d", i)))
		if err != nil {
			log.Fatal("can't send message to pull socket: %s", err.Error())
		}
		log.Println("Sent message")
		time.Sleep(time.Second)
	}
}

func subscriber(name string) {
	url := "inproc://publisher"
	sock, err := sub.NewSocket()
	if err != nil {
		log.Fatal("can't get new sub socket: %s", err.Error())
	}
	if err = sock.Dial(url); err != nil {
		log.Fatal("can't listen on pub socket: %s", err.Error())
	}
	err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		log.Fatal("cannot subscribe: %s", err.Error())
	}
	for {
		log.Println(name, " Waiting")
		msg, err := sock.Recv()
		if err != nil {
			log.Fatal("cannot receive on sub socket: %s", err.Error())
		}
		log.Println(name, " ", string(msg))
	}
}

func main() {

	go publisher()
	time.Sleep(time.Second * 1)
	go subscriber("First")
	go subscriber("Second")

	for {
		time.Sleep(time.Minute)
	}
}
