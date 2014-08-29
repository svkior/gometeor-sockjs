package stringrand

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890abcdefghijklmnopqrstuvwxyz"

const Maxlen = 10

func RandString(N int) string {
	var buf bytes.Buffer

	for j := 0; j < N; j++ {
		buf.WriteByte(chars[rand.Intn(len(chars))])
	}
	s := buf.String()
	return s
}

func Init() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println(RandString(16))
}
