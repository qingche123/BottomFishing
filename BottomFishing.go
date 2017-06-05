package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"strings"
)

type stPing struct {
	Ping int64 `json:"ping"`
}

type stPong struct {
	Pong int64 `json:"pong"`
}

var origin = "http://127.0.0.1:7321/"
var url = "wss://api.huobi.com/ws"

func packHb(hb []byte) []byte {
	pPing := &stPing{}
	pPong := &stPong{}
	err := json.Unmarshal(hb, pPing)
	if err != nil {
		fmt.Println(err.Error())
	}

	pPong.Pong = pPing.Ping

	respPong, err := json.Marshal(pPong)
	return respPong

}

func main() {
	/*
		ltcCny1day := `{
			"sub": "market.ltccny.kline.1day",
			"id": "id1"
		}`
	*/
	ltcCnyTrade := `{
		"sub": "market.ltccny.trade.detail",
		"id": "id1"
	}`

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	order := []byte(ltcCnyTrade)
	_, err = ws.Write(order)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", order)

	for {
		var msg = make([]byte, 512)
		m, err := ws.Read(msg)
		if err != nil {
			log.Fatal(err)
		}

		//deCompress data
		b := bytes.NewBuffer(msg[0:m])
		r, _ := gzip.NewReader(b)
		unData, _ := ioutil.ReadAll(r)
		r.Close()

		//data process
		if strings.Contains(string(unData), "ping") {
			respPong := packHb(unData)
			_, err = ws.Write(respPong)
			if err != nil {
				log.Fatal(err)
			}
			continue
		}

		fmt.Print("size: ", len(unData))
		fmt.Printf(" Receive: %s\n", unData)

	}
	ws.Close()
}
