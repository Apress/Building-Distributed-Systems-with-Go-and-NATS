package main

import (
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/rep"
	_ "go.nanomsg.org/mangos/v3/transport/tcp"
	"log"
)

func main() {
	url := "tcp://127.0.0.1:40899"
	sock, err := rep.NewSocket()
	if err != nil {
		log.Fatal("can't get new rep socket: %s", err.Error())
	}
	if err = sock.Listen(url); err != nil {
		log.Fatal("can't listen on rep socket: %s", err.Error())
	}
	for {
		ctx, err := sock.OpenContext()
		if err != nil {
			log.Fatal("cannot open context: %s", err.Error())
		}
		msg, err := ctx.Recv()
		if err != nil {
			log.Fatal("cannot receive on rep socket: %s", err.Error())
		}
		log.Println(string(msg))
		go func(ctx mangos.Context) {
			ctx.Send([]byte("Received: " + string(msg)))
		}(ctx)

	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
