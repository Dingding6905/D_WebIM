package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"

	"fmt"
	"samples/WebIM/models"
)

type WebSocketController struct {
	baseController
}

func (this *WebSocketController) Get() {
	fmt.Println("=====websocket (this *WebSocketController) Get()")
	uname := this.GetString("uname")
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	this.TplName = "websocket.html"
	this.Data["IsWebSocket"] = true
	this.Data["UserName"] = uname
}

func (this *WebSocketController) Join() {
	fmt.Println("=====websocket (this *WebSocketController) Join()")
	uname := this.GetString("uname")
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	Join(uname, ws)
	defer Leave(uname)

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		publish <- newEvent(models.EVENT_MESSAGE, uname, string(p))
	}
}

func broadcastWebSocket(event models.Event) {
	fmt.Println("=====websocket broadcastWebSocket(event models.Event)")
	data, err := json.Marshal(event)
	fmt.Printf("=====websocket broadcastWebSocket json data=%s\n", data)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}

	// 判断除了第一个以外还有几个链接，给其余的都发送消息
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			fmt.Println("=====websocket subscribers=", sub.Value)
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				unsubscribe <- sub.Value.(Subscriber).Name
			}
		}
	}
}
