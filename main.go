package main

import (
	"fmt"
	"github.com/hibiken/asynq"
	"gomind/tasks"
	"log"
	"runtime"
)

//func worker(jobs <-chan models.Hand, results chan<- models.FlopHand, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
//	defer wg.Done()
//	for job := range jobs {
//		strength := calculateHandStrength(job.HoleCards, job.CommunityCards)
//		results <- models.FlopHand{
//			Hand: strings.Join(job.HoleCards, "") + strings.Join(
//				job.CommunityCards,
//				"",
//			), HoleCards: strings.Join(job.HoleCards, ""),
//			Board: strings.Join(
//				job.HoleCards,
//				"",
//			),
//			Wins: strength.Wins,
//			Loss: strength.Loss, Draws: strength.Draws,
//		}
//		bar.Add(1)
//	}
//}

const redisAddr = "127.0.0.1:6379"

func enqueuTask() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------

	task, err := tasks.NewEmailDeliveryTask(42, "some:template:id")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

}

func main() {
	fmt.Println(runtime.NumCPU())
	//dsn := "host=localhost user=postgres password=example  port=5555"
	//db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//err := db.AutoMigrate(&models.FlopHand{})
	//if err != nil {
	//	return
	//}
	//hands := generateHandCombinations(3)
	//
	//bar := progressbar.Default(int64(len(hands)))
	//
	//const numWorkers = 8
	//
	//jobs := make(chan models.Hand, len(hands))
	//results := make(chan models.FlopHand, len(hands))
	//var wg sync.WaitGroup
	//
	//// Worker starten
	//for i := 0; i < numWorkers; i++ {
	//	wg.Add(1)
	//	go worker(jobs, results, &wg, bar)
	//}
	//
	//// Jobs in den Kanal senden
	//for _, item := range hands {
	//	jobs <- item
	//
	//}
	//close(jobs) // Keine weiteren Jobs mehr
	//
	//wg.Wait()
	//
	//toSave := make([]models.FlopHand, 0)
	//for result := range results {
	//	toSave = append(toSave, result)
	//}
	//
	//db.CreateInBatches(toSave, 100)
	//file, _ := os.Create("flopAbstraction.gob")
	//
	//defer file.Close()
	//
	//encoder := gob.NewEncoder(file)
	//encoder.Encode(toSave)

}
