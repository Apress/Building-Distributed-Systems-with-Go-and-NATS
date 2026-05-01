package main

import (
	"fmt"
	"go.nanomsg.org/mangos/v3/protocol/pub"
	_ "go.nanomsg.org/mangos/v3/transport/ipc"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	ps := os.Getenv("PORTS")
	pp := strings.Split(ps, ",")

	sock, err := pub.NewSocket()
	if err != nil {
		log.Fatal("can't get new push socket: %s", err.Error())
	}
	for _, port := range pp {
		url := fmt.Sprintf("ipc://%s", port)
		if err = sock.Dial(url); err != nil {
			log.Fatal("can't listen on push socket: %s - URL: %s", err.Error(), url)
		}
		log.Println("Dialling ", url)
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
