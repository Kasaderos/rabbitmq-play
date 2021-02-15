package main

import "log"

func Case1() {
	sub, err := NewRMQReceiver()
	if err != nil {
		log.Println(err)
		return
	}
	defer sub.Close()
	flag := false
	log.Println("case1 started")
	for msg := range sub.DeliveryChan {
		if flag {
			err := msg.Nack(false, false)
			if err != nil {
				log.Println("[nack] failed", err)
			}
			flag = false
			log.Println("[nack] success", string(msg.Body))
		} else {
			err := msg.Ack(true)
			if err != nil {
				log.Println("[ack] failed", err)
			}
			flag = true
			log.Println("[ack] success", string(msg.Body))
		}
	}
}

func main() {
	Case1()
}
