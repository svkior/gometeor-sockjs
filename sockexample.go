package main

import (
	"github.com/igm/pubsub"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
	"net/http"
)

var chat pubsub.Publisher

func main() {
	http.Handle("/echo/", sockjs.NewHandler("/echo", sockjs.DefaultOptions, echoHandler))
	http.Handle("/", http.FileServer(http.Dir("web/")))
	log.Println("Server started on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func echoHandler(session sockjs.Session) {
	log.Println("new sockjs session estabilished")
	var closedSession = make(chan struct{})
	chat.Publish("[info] new participant joined chat")
	defer chat.Publish("[info] participant left chat")
	go func() {
		reader, _ := chat.SubChannel(nil)
		for {
			select {
			case <-closedSession:
				log.Println("Connection closed!")
				return
			case msg := <-reader:
				if err := session.Send(msg.(string)); err != nil {
					return
				}
			}
		}
	}()
	for {
		msg, err := session.Recv()
		if err == nil {
			chat.Publish(msg)
			continue
		}
		log.Println(err.Error())
		break
	}
	close(closedSession)
	log.Println("sockjs session closed")
}
