package main

import (
	"go.nanomsg.org/mangos/v3/protocol/pull"
	_ "go.nanomsg.org/mangos/v3/transport/tcp"
	"log"
)

func main() {
	//TIP Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined or highlighted text
	// to see how GoLand suggests fixing it.
	url := "tcp://127.0.0.1:40899"
	sock, err := pull.NewSocket()
	if err != nil {
		log.Fatal("can't get new rep socket: %s", err.Error())
	}
	if err = sock.Listen(url); err != nil {
		log.Fatal("can't listen on pull socket: %s", err.Error())
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

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
