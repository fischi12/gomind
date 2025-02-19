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

func batchSlice(data []services.Hand, batchSize int) [][]services.Hand {
	var batches [][]services.Hand
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data) // Verhindert Out-of-Bounds-Fehler
		}
		batches = append(batches, data[i:end])
	}
	return batches
}

const redisAddr = "192.168.2.188:6379"

func enqueuTask() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------
	hands := services.GenerateHandCombinations(3)
	batches := batchSlice(hands, 100)
	for _, batch := range batches {
		task, err := tasks.NewHandAbstractionFlopTask(batch)
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
