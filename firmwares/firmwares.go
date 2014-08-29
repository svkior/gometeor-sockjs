package firmwares

import (
	"../stringrand"
	"fmt"
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
	firmwares []firmware
}

func (fw *Firmwares) Add(f firmware) {
	f.id = stringrand.RandString(16)
	//TODO: need to find duplications in random generation
	fw.firmwares = append(fw.firmwares, f)
}

func TestInitFirmwares(fw *Firmwares) {
	fw.Add(firmware{
		url:         "http://www.ya.ru",
		fwname:      "Хреновая прошивка",
		description: "Вот такая прошивка",
		author:      "Sergey V. Kior",
	})
}

func (fw Firmwares) GetAllJSON() (s chan string) {
	s = make(chan string)
	go func() {
		for _, v := range fw.firmwares {
			s <- fmt.Sprintf(
				"{\"msg\": \"added\", \"collection\":\"firmwares\", \"id\": \"%s\", \"fields\":{\"url\":\"%s\",\"fwname\":\"%s\",\"description\":\"%s\",\"author\":\"%s\",\"downloaded\": %t }}",
				v.id,
				v.url,
				v.fwname,
				v.description,
				v.author,
				v.downloaded,
			)
		}
		close(s)
	}()
	return
}
