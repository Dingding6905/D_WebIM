package controllers

import (
	"fmt"
	"samples/WebIM/models"
)

type LongPollingController struct {
	baseController
}

func (this *LongPollingController) Join() {
	fmt.Println("=====longpolling (this *LongPollingController) Join()")
	uname := this.GetString("uname")
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	Join(uname, nil)

	this.TplName = "longpolling.html"
	this.Data["IsLongPolling"] = true
	this.Data["UserName"] = uname
}

func (this *LongPollingController) Post() {
	fmt.Println("=====longpolling (this *LongPollingController) Post()")
	this.TplName = "longpolling.html"

	uname := this.GetString("uname")
	content := this.GetString("content")
	if len(uname) == 0 || len(content) == 0 {
		return
	}

	fmt.Println("=====longpolling Post()=", newEvent(models.EVENT_MESSAGE, uname, content))
	publish <- newEvent(models.EVENT_MESSAGE, uname, content)
}

// 此函数执行完成和 longpolling.js 中 getJSON 第二个函数进行交互
func (this *LongPollingController) Fetch() {
	fmt.Println("=====longpolling (this *LongPollingController) Fetch()")
	lastReceived, err := this.GetInt("lastReceived")
	if err != nil {
		return
	}

	events := models.GetEvents(int(lastReceived))
	if len(events) > 0 {
		this.Data["json"] = events
		this.ServeJSON()
		return
	}

	ch := make(chan bool)
	waitingList.PushBack(ch)

	/*****************/
	for e := waitingList.Front(); e != nil; e = e.Next() {
		fmt.Println("=====longpooling Fetch waitingList=", e.Value)
	}
	/****************/

	<-ch

	this.Data["json"] = models.GetEvents(int(lastReceived))
	this.ServeJSON()
}
