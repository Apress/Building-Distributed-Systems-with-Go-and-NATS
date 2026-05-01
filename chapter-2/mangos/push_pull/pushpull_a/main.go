package main

import (
	"fmt"
	"go.nanomsg.org/mangos/v3/protocol/pull"
	"go.nanomsg.org/mangos/v3/protocol/push"
	_ "go.nanomsg.org/mangos/v3/transport/tcp"
	"log"
)

func receiver(ch chan bool) {
	plurl := "tcp://127.0.0.1:40899"
	plsock, err := pull.NewSocket()
	if err != nil {
		log.Fatal("can't get new rep socket: %s", err.Error())
	}
	if err = plsock.Dial(plurl); err != nil {
		log.Fatalf("can't listen on push socket: %s URL: %s", err.Error(), plurl)
	}

	for {
		msg, err := plsock.Recv()
		if err != nil {
			log.Fatal("cannot send to pull socket: %s", err.Error())
		}
		log.Println(string(msg))

		ch <- true
	}

}
func main() {

	ch := make(chan bool)

	psurl := "tcp://127.0.0.1:40898"
	pssock, err := push.NewSocket()
	if err != nil {
		log.Fatal("can't get new rep socket: %s", err.Error())
	}
	if err = pssock.Dial(psurl); err != nil {
		log.Fatalf("can't listen on pull socket: %s URL: %s", err.Error(), psurl)
	}

	go receiver(ch)
	i := 0
	for {
		i += 1
		err := pssock.Send([]byte(fmt.Sprintf("message %d", i)))
		if err != nil {
			log.Fatal("cannot send to pull socket: %s", err.Error())
		}
		log.Println("Sedning!")

		if i == 100 {
			break
		}
	}

	for i := 0; i < 10; i++ {
		<-ch
	}
	log.Println("Finished")
}
