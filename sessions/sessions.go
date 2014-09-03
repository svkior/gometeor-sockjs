package sessions

import (
	"../collections"
	"../stringrand"
	"./meteorMethods"
	"./meteorSession"
	"encoding/json"
	"fmt"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
	"reflect"
)

type metCol struct {
	colName   string // название колекции
	colStruct *collections.Collectioner
}

// Коллекция сессий с метеором
type MeteorSessions struct {
	sessions    []meteorSession.MeteorSession
	collections []metCol
	methods     []meteorMethods.MeteorMethod
}

func (ms *MeteorSessions) Subscribe(idx int, name string) {
	// ms.Subscribe(ccId, subName)
	//fmt.Println("Is Subscribed: ", s.IsUbscribed(name))
	ms.sessions[idx].Subscribe(name)
}

func (ms *MeteorSessions) GetAllJSON(name string) chan string {
	return (*ms.GetCollection(name)).GetAllJSON()
}

func (ms *MeteorSessions) SubscribeChan(name string) chan string {
	return (*ms.GetCollection(name)).SubscribeChan()
}

func (ms *MeteorSessions) GetSessionId(idx int) string {
	return ms.sessions[idx].GetId()
}

func (ms *MeteorSessions) GetSessionIdx(hash string) int {

	for idx := 0; idx < len(ms.sessions); idx++ {
		if ms.sessions[idx].GetId() == hash {
			return idx
		}
	}
	return -1
}

func (ms *MeteorSessions) IsSubscribed(idx int, name string) bool {
	return ms.sessions[idx].IsSubcribed(name)
}

func (ms *MeteorSessions) AddCollection(colName string, collection collections.Collectioner) {

	ms.collections = append(ms.collections, metCol{colName: colName, colStruct: &collection})
	ms.AddMethod("/"+colName+"/insert", collection.Insert)
}

func (ms *MeteorSessions) AddMethod(mName string, f func(params interface{}) string) {
	m := meteorMethods.Create(mName, f)
	ms.methods = append(ms.methods, m)
}

func (ms *MeteorSessions) GetMethodIdx(name string) int {
	for i := 0; i < len(ms.methods); i++ {
		if ms.methods[i].NameEquals(name) {
			return i
		}
	}
	return -1
}

func (ms *MeteorSessions) CallMethodByIdx(idx int, params interface{}) string {
	return ms.methods[idx].CallMethod(params)
}

// Добавить сессию в коллекцию
func (m *MeteorSessions) Append(session string) {
	m.sessions = append(m.sessions, meteorSession.Create(session))
}

func (m *MeteorSessions) Length() int {
	return len(m.sessions)
}

func (ms *MeteorSessions) MeteorHandler(session sockjs.Session) {
	log.Println("new sockjs session estabilished")
	//fmt.Println("Number of current sessions:", ms.Length())

	// Канал, в который будут сыпаться все
	transmitter := make(chan string)

	var ccId int
	//	var ms meteorSession
	//var closedSession = make(chan struct{})

	go func() {
		for {
			select {
			//case <-closedSession:
			//				log.Println("Connection closed!")
			//				return
			case msg := <-transmitter:
				if err := session.Send(msg); err != nil {
					//				if err := session.Send(msg.(string)); err != nil {
					return
				}
			}
		}
	}()

	for {
		msg, err := session.Recv()
		if err == nil {
			//fmt.Println("Received message: ", msg)
			var f interface{}
			err = json.Unmarshal([]byte(msg), &f)
			if err != nil {
				fmt.Println("Error: ", err.Error())
			} else {
				//log.Println("Got message: ", f)

				m := f.(map[string]interface{})

				var neeedPrint bool

				switch m["msg"] {
				case "connect": /// <<<<<<< CONNECT
					//fmt.Println("Connect from client")
					//
					//  m["session"] - random session id
					//  m["version"] - proposed version to connect
					//
					//fmt.Println("Request Version: ", m["version"])
					if m["version"] != "pre2" {
						// msg: failed
						// version: "pre2"
						transmitter <- "{\"msg\" :  \"failed\", \"version\" : \"pre2\" }"
					}

					if m["session"] == nil {
						sessionId := stringrand.RandString(16)
						log.Println("New connection ID: ", sessionId)
						ms.Append(sessionId)

						transmitter <- "{\"msg\" :  \"connected\", \"session\" : \"" + sessionId + "\" }"

						//						if err := session.Send("{\"msg\" :  \"connected\", \"session\" : \"" + sessionId + "\" }"); err != nil {
						//							return
						//						}

						ccId = ms.GetSessionIdx(sessionId)
						if ccId == -1 {
							log.Panic("EEEEE: FFFUUUUKKKKK")
						}

					} else {
						var sess = m["session"].(string)

						//fmt.Println("Request to reconnect to old session: ", sess)
						ccId = ms.GetSessionIdx(sess)
						if ccId == -1 {
							fmt.Println("There is no old session, creating new one")
							sessionId := stringrand.RandString(16)
							log.Println("New connection ID: ", sessionId)
							ms.Append(sessionId)
							transmitter <- "{\"msg\" :  \"connected\", \"session\" : \"" + sessionId + "\" }"
							ccId = ms.GetSessionIdx(sessionId)
						} else {
							transmitter <- "{\"msg\" :  \"connected\", \"session\" : \"" + ms.GetSessionId(ccId) + "\" }"
						}

					}
					/// **** PING
				case "ping":
					//fmt.Println("Ping (Unrealized)")
					if m["id"] == nil {
						transmitter <- "{\"msg\" :  \"pong\"}"
					} else {
						var idString string
						idString = m["id"].(string)

						transmitter <- "{\"msg\" :  \"pong\", \"id\" : \"" + idString + "\"}"
					}
					//fmt.Println("Ping Packet:", m)
					//neeedPrint = true

				case "sub":
					//fmt.Println("=> Subscribe (Unrealized)")
					subName := m["name"].(string)

					if ms.HasCollection(subName) {
						//log.Println("We have collection ", subName)
						if ms.IsSubscribed(ccId, subName) == false {
							log.Println("Subscribing to:", subName)
							ms.Subscribe(ccId, subName)
							// Отправляем все записи согласно подписке
							// NB!! Это корректное использование range!!!!!
							for msg := range ms.GetAllJSON(subName) {
								//fmt.Println("GotMsg:", msg)
								transmitter <- msg
							}

							// Подписываемся на коллекции
							go func() {
								var firmchan chan string
								// FIXME: Нужно придумать как грохать эту хрень когда клиент закрывается
								firmchan = ms.SubscribeChan(subName)
								for {
									select {
									case msg2 := <-firmchan:
										log.Println("Send changes from collection to client: ", msg2)
										transmitter <- msg2
									}
								}
							}()

						} else {
							log.Println("Already subscribed to:", subName)
						}

						// Говорим, что все изменения вычитаны
						transmitter <- "{\"msg\": \"ready\", \"subs\": [\"" + m["id"].(string) + "\"]}"

					} else {
						log.Println("We don't provide subscription to : ", subName)
						// There is no Subscription to autoupdate
						transmitter <- "{\"msg\" :  \"nosub\", \"id\" : \"" + subName + "\", \"error\" : \"sub-not-found\" }"
					}

					//neeedPrint = true
				case "unsub":
					fmt.Println("Unsubscribe (Unrealized)")

				// METHOD
				case "method":
					//fmt.Println("Method call (Unrealized)")
					m2 := m["params"].([]interface{})
					var methodId string
					var methodName string

					methodId = m["id"].(string)
					methodName = m["method"].(string)

					fmt.Println("Method Name: ", methodName)
					fmt.Println("Method ID: ", methodId)
					fmt.Println("Params: ", m2)

					if midx := ms.GetMethodIdx(methodName); midx == -1 {
						log.Println("There is no such method", methodName)
						transmitter <- "{\"msg\" :  \"result\", \"id\" : \"" + methodId + "\", \"error\" : \"method-not-found\" }"
					} else {
						log.Println("We have method:", methodName)
						go func() {
							result := ms.CallMethodByIdx(midx, m2)
							fmt.Println("REEEEZUUULT: ", result)
							transmitter <- "{\"msg\" :  \"result\", \"id\" : \"" + methodId + "\", \"result\" : \"+ result +\" }"
							transmitter <- "{\"msg\" :  \"updated\", \"methods\" : [\"" + methodId + "\"] }"
						}()
					}

					neeedPrint = false
				default:
					fmt.Println("Unknown Message type: ", m["msg"])
					neeedPrint = true
				}
				if neeedPrint {
					for k, v := range m {
						fmt.Println(k, " => ", v, " <Type> => ", reflect.TypeOf(v))
					}
				}

			}
			continue
		}
		log.Println(err.Error())
		break
	}
	//close(closedSession)
	log.Println("sockjs session closed")

}

func (ms *MeteorSessions) HasCollection(colName string) bool {
	for i := 0; i < len(ms.collections); i++ {
		if ms.collections[i].colName == colName {
			return true
		}
	}
	return false
}

func (ms *MeteorSessions) GetCollection(colName string) (col *collections.Collectioner) {
	for i := 0; i < len(ms.collections); i++ {
		if ms.collections[i].colName == colName {
			return ms.collections[i].colStruct
		}
	}
	return nil
}
