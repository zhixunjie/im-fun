package gmodel

import (
	"fmt"
	"log"
	"testing"
)

func TestIdSession(t *testing.T) {
	id1 := NewUserComponentId(1001)
	id2 := NewUserComponentId(1002)
	id3 := NewGroupComponentId(10)

	fmt.Println("单聊", NewSessionId(id1, id2))
	fmt.Println("单聊", NewSessionId(id2, id1))
	fmt.Println("群聊", NewSessionId(id1, id3))
	fmt.Println("群聊", NewSessionId(id3, id1))
	fmt.Println("群聊", NewSessionId(id2, id3))
	fmt.Println("群聊", NewSessionId(id3, id2))
}

func TestSort(t *testing.T) {
	id1 := NewComponentId(1005, 2)
	id2 := NewComponentId(1004, 1)
	fmt.Println(id1.Sort(id2))
}

func TestParseSessionId(t *testing.T) {
	sessionId := NewSessionId(NewUserComponentId(1001), NewGroupComponentId(100000000001))
	result, err := sessionId.Parse()
	if err != nil {
		log.Fatalln()
	}
	fmt.Println(result, result.Ids[0])

	sessionId = NewSessionId(NewUserComponentId(1001), NewUserComponentId(1002))
	result, err = sessionId.Parse()
	if err != nil {
		log.Fatalln()
	}
	fmt.Println(result, result.Ids[0], result.Ids[1])

	sessionId = NewSessionId(NewUserComponentId(1001), NewRobotComponentId(111111))
	result, err = sessionId.Parse()
	if err != nil {
		log.Fatalln()
	}
	fmt.Println(result, result.Ids[0], result.Ids[1])
}
