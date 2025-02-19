package main

import (
	"fmt"
	"github.com/hibiken/asynq"
	"gomind/database"
	"gomind/models"
	"gomind/services"
	"gomind/tasks"
	"log"
	"runtime"
)

const redisAddr = "127.0.0.1:6379"

func enqueuTask() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------
	hands := services.GenerateHandCombinations(3)
	for _, hand := range hands {
		task, err := tasks.NewHandAbstractionFlopTask(hand)
		if err != nil {
			log.Fatalf("could not create task: %v", err)
		}
		info, err := client.Enqueue(task)
		if err != nil {
			log.Fatalf("could not enqueue task: %v", err)
		}
		log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	}

}

func main() {
	fmt.Println(runtime.NumCPU())

	db := database.GetDB()
	err := db.AutoMigrate(&models.FlopHand{})
	if err != nil {
		return
	}

	enqueuTask()

}
