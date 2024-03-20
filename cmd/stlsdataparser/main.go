package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/config"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/parser"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/processor"
	"github.com/jacksonbarreto/WebGateScanner-kafka/producer"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	const configFilePath = ""
	var (
		processing = make(map[string]bool)
		lock       = sync.Mutex{}
	)
	config.InitConfig(configFilePath)
	errorPath := config.App().ErrorParsePath
	pathToWatch := config.App().PathToWatch
	totalWorkers := config.App().Workers
	kafkaProducer, producerErr := producer.NewProducer(config.Kafka().TopicsProducer[0], config.Kafka().Brokers, config.Kafka().MaxRetry)
	if producerErr != nil {
		panic(producerErr)
	}
	defer kafkaProducer.Close()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if _, err := os.Stat(errorPath); os.IsNotExist(err) {
		os.Mkdir(errorPath, os.ModePerm)
	}

	if _, err := os.Stat(pathToWatch); os.IsNotExist(err) {
		os.Mkdir(pathToWatch, os.ModePerm)
	}

	filesToProcess := make(chan string, 100)

	for i := 0; i < totalWorkers; i++ {
		worker := processor.NewWorker(kafkaProducer, errorPath, &lock, parser.NewParser(), processing)
		go worker.Do(filesToProcess)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create && strings.HasSuffix(event.Name, ".done") {
					log.Println("NewProcessor file detected:", event.Name)
					originalFileName := strings.TrimSuffix(event.Name, ".done")
					lock.Lock()
					if !processing[event.Name] {
						processing[event.Name] = true
						filesToProcess <- originalFileName
					}
					lock.Unlock()
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(pathToWatch)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
