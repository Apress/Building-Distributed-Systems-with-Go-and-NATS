package main

import (
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/sub"
	_ "go.nanomsg.org/mangos/v3/transport/tcp"
	"log"
	"os"
	"strings"
)

func main() {
	ts := os.Getenv("TOPICS")
	tt := strings.Split(ts, ",")
	url := "tcp://127.0.0.1:40899"
	sock, err := sub.NewSocket()
	if err != nil {
		log.Fatal("can't get new rep socket: %s", err.Error())
	}
	if err = sock.Dial(url); err != nil {
		log.Fatal("can't listen on pull socket: %s", err.Error())
	}
	for _, t := range tt {
		err = sock.SetOption(mangos.OptionSubscribe, []byte(t))
		if err != nil {
			log.Fatal("cannot subscribe: %s", err.Error())
		}
	}

	for {
		log.Println("Waiting")
		msg, err := sock.Recv()
		if err != nil {
			log.Fatal("cannot receive on pull socket: %s", err.Error())
		}
		log.Println(string(msg))
	}
}
