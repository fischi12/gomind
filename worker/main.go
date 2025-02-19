package main

import (
	"fmt"
	"github.com/hibiken/asynq"
	"gomind/tasks"
	"log"
	"runtime"
)

const redisAddr = "192.168.2.188:6379"

func main() {
	fmt.Println("Worker")
	fmt.Println(runtime.NumCPU())

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: runtime.NumCPU() - 1,
		},
	)

	taskHandler := tasks.NewTaskHandler()

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeHandAbstractionFlop, taskHandler.HandleHandAbstractionFlopTask)
	mux.HandleFunc(tasks.TypeHandAbstractionTurn, taskHandler.HandleHandAbstractionTurnTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
