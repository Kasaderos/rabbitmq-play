package main

import (
	"log"
	"time"

	"github.com/google/uuid"
)

func main() {
	pub, err := NewRMQSender()
	if err != nil {
		log.Println(err)
		return
	}
	defer pub.Close()

	// case 1
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		msg := uuid.New().String()
		err := pub.Send([]byte(msg))
		if err != nil {
			log.Println("[send] failed", err)
		} else {
			log.Println("[pub] sent message", msg)
		}
	}
}
