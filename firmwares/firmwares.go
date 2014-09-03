package firmwares

import (
	//"../mydebug"
	"../stringrand"
	"code.google.com/p/go.net/html"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"reflect"
	"encoding/json"
	"strings"
)

type firmware struct {
	id          string // ID записи (эмуляция ID в MongoDB)
	url         string // Ссылка на прошивку
	fwname      string // имя прошивки
	description string // Описание прошивки от Макса
	author      string // Автор прошивки
	downloaded  bool   // Скачивалась ли прошивка
}

type Firmwares struct {
	firmwares   []firmware
	subscribers []chan string
}

func (fw *Firmwares) Remove(params interface{}) string {
	/* From Meteor DDP Analyzer
	2  OUT  3646  {
		"msg":"method",
		"method":"/firmwares/remove",
		"params":[{
			"_id":"cMFtnvjD6TaLFZQkH"
		}],
		"id":"1"}
	2  IN   6  {
		"msg":"removed",
		"collection":"firmwares",
		"id":"cMFtnvjD6TaLFZQkH"
	}
	2  IN   2  {
		"msg":"result",
		"id":"1",
		"result":1
	}
	2  IN   1  {
		"msg":"updated",
		"methods":["1"]
	}
	*/

	// TODO: Remove record by ID
	// TODO: Format result
	return "[]"
}

func (fw *Firmwares) Insert(params interface{}) string {
	/* From Meteor DDP Analyzer
	1  OUT  13271  {
		"msg":"method",
	   	"method":"/firmwares/insert",
	   	"params":[{
	   		"author":"svkior",
	   		"url":"http://localhost:3000/superproshivha1.bit",
	   		"fwname":"top_arm_from_hell.bit",
	   		"description":"qwe"
	   	}],
	   	"id":"1",
	   	"randomSeed":"38842ec8265a97554324"
	}
	1  IN   21  {
		"msg":"result",
	   	"id":"1",
	   	"result":[{
	   		"author":"svkior",
	   		"url":"http://localhost:3000/superproshivha1.bit",
	   		"fwname":"top_arm_from_hell.bit",
	   		"description":"qwe",
	   		"_id":"5oGxoS5FzLb9u6vrQ"}
	   ]
	}
	1  IN   0  {
		"msg":"added",
		"collection":"firmwares",
		"id":"5oGxoS5FzLb9u6vrQ",
		"fields":{
			"author":"svkior",
			"url":"http://localhost:3000/superproshivha1.bit",
			"fwname":"top_arm_from_hell.bit",
			"description":"qwe"
		}
	}
	1  IN   2  {
		"msg":"updated",
		"methods":["1"]
	}
	*/

	fmt.Println("NEEED TO INSERT DOCUMENT: ", params)

	m2 := params.([]interface{})
	methodParams := m2[0].(map[string]interface{})

	fmt.Println("Method params: ", methodParams)

	if methodParams["fwname"] == nil {
		methodParams["fwname"] = "NoName_" + stringrand.RandString(4)
	}

	fwname := methodParams["fwname"].(string)
	for par, val := range methodParams {
		if par != "fwname" {
			fmt.Println("Update ", par, " => ", val)
			fw.UpdateFirmwareInfoByName(fwname, par, val.(string))
		}
	}
	methodParams["_id"] = fw.GetFwByName(fwname).id

	marshalled, _ := json.Marshal(methodParams)

	return "[" + string(marshalled) + "]"
}

func (fw *Firmwares) SubscribeChan() (s chan string) {
	s = make(chan string)
	fw.subscribers = append(fw.subscribers, s)
	return
}

func (fw *Firmwares) PushChanges(id string, added bool) {
	for i := 0; i < len(fw.firmwares); i++ {
		if fw.firmwares[i].id == id {
			//TODO: Real changes
			fmt.Println("Changed: ", fw.firmwares[i])
			for j := 0; j < len(fw.subscribers); j++ {
				if added {
					fw.subscribers[j] <- fw.GetAddedMsgByIdx(i)
				} else {
					fw.subscribers[j] <- fw.GetChangedMsgByIdx(i)
				}
			}
		}
	}
}

func (fw *Firmwares) GetFwByName(fwname string) *firmware {
	for i := 0; i < len(fw.firmwares); i++ {
		//fmt.Println("GetFwByName: ", fw.firmwares[i])
		if fw.firmwares[i].fwname == fwname {
			//fmt.Println("Found: ", fwname)
			return &fw.firmwares[i]
		}
	}
	return nil
}

func (fw *Firmwares) UpdateFirmwareInfoByName(fwname string, param string, value string) {
	var added bool
	fmt.Printf("UpdateFirmwareInfoByName %p\n", fw)

	f := fw.GetFwByName(fwname)
	if f == nil {
		fw.Add(firmware{id: stringrand.RandString(16), fwname: fwname})
		f = fw.GetFwByName(fwname)
		added = true
	}
	switch param {
	case "url":
		f.url = value
	case "description":
		f.description = value
	}
	fw.PushChanges(f.id, added)

	//fmt.Println("Finish UpdateFirmwareInfoByName")
}

func (fw *Firmwares) Add(f firmware) {
	f.id = stringrand.RandString(16)
	//TODO: need to find duplications in random generation
	fw.firmwares = append(fw.firmwares, f)
	//TODO: Send info about add
}

func TestInitFirmwares(fw *Firmwares) {
	fw.Add(firmware{
		url:         "http://www.ya.ru",
		fwname:      "Хреновая прошивка",
		description: "Вот такая прошивка",
		author:      "Sergey V. Kior",
	})
}

func (fw *Firmwares) GetAddedMsgByIdx(i int) string {
	msg := fmt.Sprintf(
		"{\"msg\": \"added\", \"collection\":\"firmwares\", \"id\": \"%s\", \"fields\":{\"url\":\"%s\",\"fwname\":\"%s\",\"description\":\"%s\",\"author\":\"%s\",\"downloaded\": %t }}",
		fw.firmwares[i].id,
		fw.firmwares[i].url,
		fw.firmwares[i].fwname,
		fw.firmwares[i].description,
		fw.firmwares[i].author,
		fw.firmwares[i].downloaded,
	)
	return msg
}

func (fw *Firmwares) GetChangedMsgByIdx(i int) string {
	msg := fmt.Sprintf(
		"{\"msg\": \"changed\", \"collection\":\"firmwares\", \"id\": \"%s\", \"fields\":{\"url\":\"%s\",\"fwname\":\"%s\",\"description\":\"%s\",\"author\":\"%s\",\"downloaded\": %t }}",
		fw.firmwares[i].id,
		fw.firmwares[i].url,
		fw.firmwares[i].fwname,
		fw.firmwares[i].description,
		fw.firmwares[i].author,
		fw.firmwares[i].downloaded,
	)
	return msg
}

func (fw *Firmwares) GetAllJSON() (s chan string) {
	s = make(chan string)
	fmt.Printf("GetAllJSON %p\n", fw)
	go func() {

		for i := 0; i < len(fw.firmwares); i++ {
			s <- fw.GetAddedMsgByIdx(i)
		}
		close(s)
	}()
	return
}

func (fw *Firmwares) Scan4DAV(params interface{}) string {
	log.Println("Scanning web directory for firmwares")
	//mydebug.PrintDebug("FW After new", fw)

	m2 := params.([]interface{})
	methodParams := m2[0].(map[string]interface{})
	var pattern string
	var url string
	for a, b := range methodParams {
		//fmt.Println("Params: ", a, " => ", b, " <Type> => ", reflect.TypeOf(b))
		switch a {
		case "pattern":
			pattern = b.(string)
		case "dirname":
			url = b.(string)
		}
	}

	// Здесь нужно запустить http клиента

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("svkior", "forserveryf[")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
	defer resp.Body.Close()
	/*
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	*/

	doc, err := html.Parse(resp.Body)

	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attrs := range n.Attr {
				if attrs.Key == "href" {
					if strings.Contains(attrs.Val, pattern) {
						id := ""
						if strings.Contains(attrs.Val, ".bit") {
							id = strings.Replace(attrs.Val, ".bit", "", -1)
							//fmt.Println(" BIT -> ", attrs.Val, " ID = ", id)
							fw.UpdateFirmwareInfoByName(id, "url", url+attrs.Val)
						} else if strings.Contains(attrs.Val, "_info") {
							id = strings.Replace(attrs.Val, "_info", "", -1)
							//fmt.Println(" INF -> ", attrs.Val, " ID = ", id)

							req2, err := http.NewRequest("GET", url+attrs.Val, nil)
							req2.SetBasicAuth("svkior", "forserveryf[")
							resp2, err := client.Do(req2)
							if err != nil {
								fmt.Printf("Error : %s", err)
							}
							defer resp2.Body.Close()
							body2, err := ioutil.ReadAll(resp.Body)
							//fmt.Println(body2)
							fw.UpdateFirmwareInfoByName(id, "description", string(body2))
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	fmt.Println("FINISH CALLING METHOD SCAN4DAV")
	return "[]"
}
