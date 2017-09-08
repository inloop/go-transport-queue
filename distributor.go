package main

import (
	"fmt"
	"time"

	"github.com/inloop/go-transport-queue/model"
)

type distributor struct {
	queue     *Queue
	transport model.Transport
	running   bool
}

func (d *distributor) start(rate time.Duration, batchSize int) {
	d.running = true
	go func() {
		throttle := time.Tick(rate)
		for range throttle {
			messages := []model.TransportMessage{}
			for i := 0; i < batchSize; i++ {
				if item, _ := d.queue.Pop(); item != nil {
					messages = append(messages, *item)
				} else {
					break
				}
			}
			d.sendMessages(messages, 0)
		}
	}()
}

func (d *distributor) sendMessages(messages []model.TransportMessage, iteration int) {
	if iteration > 3 {
		fmt.Println("stopped sending messages after", iteration-1, "tries")
		return
	}

	if err := d.transport.SendMessages(messages); err != nil {
		fmt.Println("failed to send messages, try:", iteration+1, "err:", err)
		time.Sleep(time.Second)
		d.sendMessages(messages, iteration+1)
	}
}

func (d *distributor) stop() {
	d.running = false
}
