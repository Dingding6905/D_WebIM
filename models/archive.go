package models

import (
	"container/list"
	"fmt"
)

type EventType int

const (
	EVENT_JOIN = iota
	EVENT_LEAVE
	EVENT_MESSAGE
)

type Event struct {
	Type      EventType
	User      string
	Timestamp int
	Content   string
}

const archiveSize = 20

// 挂 Event 用户及其详细信息结构体
var archive = list.New()

// 将 Event 结构体挂到 archive 链上
func NewArchive(event Event) {
	fmt.Println("===== archive NewArchive(event Event)")
	// 如果超长则删除第一个
	if archive.Len() >= archiveSize {
		archive.Remove(archive.Front())
	}
	archive.PushBack(event)

	/*****************/
	for e := archive.Front(); e != nil; e = e.Next() {
		fmt.Println("=====archive NewArchive archive=", e.Value)
	}
	/****************/
}

func GetEvents(lastReceived int) []Event {
	fmt.Println("=====archive GetEvents(lastReceived int) []Event")
	events := make([]Event, 0, archive.Len())
	for event := archive.Front(); event != nil; event = event.Next() {
		e := event.Value.(Event)
		if e.Timestamp > int(lastReceived) {
			// 如果 lastReceived 小于 e.Timestamp 则把 e 追加到 events 切片后
			events = append(events, e)
		}
	}
	fmt.Println("===========archive events=", events)
	return events
}
