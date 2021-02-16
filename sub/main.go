package main

import "log"

func Case() {
	sub, err := NewRMQReceiver()
	if err != nil {
		log.Println(err)
		return
	}
	defer sub.Close()
	log.Println("nack case with expiration started")
	for msg := range sub.DeliveryChan {
		err := msg.Nack(false, true)
		if err != nil {
			log.Println("[nack] failed", err)
		}
		log.Println("[nack] success", msg.RoutingKey, string(msg.Body))
	}
}

func main() {
	Case()
}
