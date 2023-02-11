package app

import (
	"log"
	"time"
)

func Run() {
	log.Println("run")
	duration := time.Minute * 20
	time.Sleep(duration)
}
