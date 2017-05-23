package controllers

import (
	"container/list"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"

	"fmt"
	"samples/WebIM/models"
)

type Subscription struct {
	Archive []models.Event
	New     <-chan models.Event
}

func newEvent(ep models.EventType, user, msg string) models.Event {
	fmt.Println("=====chatroom newEvent(ep models.EventType, user, msg string) models.Event=", models.Event{ep, user, int(time.Now().Unix()), msg})
	return models.Event{ep, user, int(time.Now().Unix()), msg}
}

func Join(user string, ws *websocket.Conn) {
	fmt.Println("=====chatroom Join(user string, ws *websocket.Conn)=", Subscriber{Name: user, Conn: ws})
	// 将用户名写入 Subscriber 结构体中并通过 channel 发送到 subscribe
	subscribe <- Subscriber{Name: user, Conn: ws}
}

func Leave(user string) {
	fmt.Println("=====chatroom Leave(user string)=", user)
	unsubscribe <- user
}

type Subscriber struct {
	Name string
	Conn *websocket.Conn
}

var (
	subscribe   = make(chan Subscriber, 10)
	unsubscribe = make(chan string, 10)
	publish     = make(chan models.Event, 10)
	// 等待发送消息的队列，和 getJSON 函数相关联阻塞
	waitingList = list.New()
	// Subscriber 用户名链表
	subscribers = list.New()
)

func chatroom() {
	fmt.Println("=====chatroom chatroom()")
	for {
		select {
		case sub := <-subscribe:
			fmt.Println("=====chatroom sub")
			// 如果用户不存在则添加 Subscriber 结构体到 subscribers 链表
			if !isUserExist(subscribers, sub.Name) {
				subscribers.PushBack(sub)
				// 将 Event 结构体返回并写入 channel
				publish <- newEvent(models.EVENT_JOIN, sub.Name, "")

				/*****************/
				for e := subscribers.Front(); e != nil; e = e.Next() {
					fmt.Println("=====chatroom subscribers=", e.Value)
				}
				/****************/

				beego.Info("New user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			} else {
				beego.Info("Old user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			}
		case event := <-publish:
			fmt.Println("=====chatroom event")
			/*****************/
			for e := waitingList.Front(); e != nil; e = e.Next() {
				fmt.Println("=====chatroom waitingList= ", e.Value)
			}
			/****************/

			for ch := waitingList.Back(); ch != nil; ch = ch.Prev() {
				ch.Value.(chan bool) <- true
				fmt.Println("=====chatroom waitingList= ", ch.Value)
				waitingList.Remove(ch)
			}

			broadcastWebSocket(event)
			models.NewArchive(event)

			if event.Type == models.EVENT_MESSAGE {
				beego.Info("Message from", event.User, ";Content:", event.Content)
			}
		case unsub := <-unsubscribe:
			fmt.Println("=====chatroom unsub")
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				fmt.Println("=====chatroom unsub subscribers=", sub.Value)
				if sub.Value.(Subscriber).Name == unsub {
					subscribers.Remove(sub)
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
						beego.Error("WebSocket closed:", unsub)
					}
					publish <- newEvent(models.EVENT_LEAVE, unsub, "")
					break
				}
			}
		}
	}
}

func init() {
	fmt.Println("=====chatroom.go init()")
	go chatroom()
}

func isUserExist(subscribers *list.List, user string) bool {
	fmt.Println("=====chatroom isUserExist(subscribers *list.List, user string) bool")
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).Name == user {
			return true
		}
	}
	return false
}
