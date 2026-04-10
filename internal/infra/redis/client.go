package redis

import "log"

func Enqueue(id string) error {
	// TODO: push to Redis queue
	log.Println("Enqueued:", id)
	return nil
}