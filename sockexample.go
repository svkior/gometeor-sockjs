package main

import (
	"./firmwares"
	"./sessions"
	"./stringrand"
	//"fmt"
	"github.com/igm/pubsub"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
	"net/http"
)

var chat pubsub.Publisher      // Из примера на pubsub
var fw firmwares.Firmwares     // База данных прошивок
var ms sessions.MeteorSessions // Сессии метеора

func main() {
	// Инициализация счетчика случайных чисел
	stringrand.Init()
	// Забиваем тестовую прошивку
	firmwares.TestInitFirmwares(&fw)

	ms.AddMethod("scan4dav", fw.Scan4DAV)
	ms.AddCollection("firmwares", fw)

	// SockJS для Meteor DDP
	http.Handle("/sockjs/", sockjs.NewHandler("/sockjs", sockjs.DefaultOptions, ms.MeteorHandler))
	// Тестовый SockJS для эхо
	http.Handle("/echo/", sockjs.NewHandler("/echo", sockjs.DefaultOptions, echoHandler))
	// Файловый сервер
	http.Handle("/", http.FileServer(http.Dir("web/")))
	// TODO: Нужно сделать темплейт, который разбирает .json от Meteor и создает главный файл html
	// Запускаем www сервер на порту 8080
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
