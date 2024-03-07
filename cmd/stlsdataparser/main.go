package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/config"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/parser"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/processor"
	"github.com/jacksonbarreto/WebGateScanner-kafka/producer"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	processing = make(map[string]bool)
	lock       = sync.Mutex{}
)

const configFilePath = ""

func main() {
	config.InitConfig(configFilePath)
	errorPath := config.App().ErrorParsePath
	pathToWatch := config.App().PathToWatch

	kafkaProducer, producerErr := producer.New(config.Kafka().TopicsProducer[0], config.Kafka().Brokers, config.Kafka().MaxRetry)
	if producerErr != nil {
		panic(producerErr)
	}
	defer kafkaProducer.Close()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Crie a pasta de erros se ela não existir
	if _, err := os.Stat(errorPath); os.IsNotExist(err) {
		os.Mkdir(errorPath, os.ModePerm)
	}

	// Worker pool para processar arquivos em paralelo
	filesToProcess := make(chan string, 100) // Buffer pode ser ajustado conforme necessário

	// Iniciando workers
	for i := 0; i < 10; i++ { // Número de workers
		go worker(filesToProcess, kafkaProducer)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("New file detected:", event.Name)
					lock.Lock()
					if !processing[event.Name] {
						processing[event.Name] = true
						filesToProcess <- event.Name
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

	// Bloqueia o main indefinidamente
	select {}
}

func worker(files <-chan string, kafkaProducer *producer.Producer) {
	for filePath := range files {
		log.Println("Processing file:", filePath)
		process := processor.New(kafkaProducer, parser.New())
		if err := process.ProcessFile(filePath); err != nil {
			log.Println("Failed to process file:", err)
			os.Rename(filePath, filepath.Join(config.App().ErrorParsePath, filepath.Base(filePath)))
		} else {
			os.Remove(filePath)
		}
		lock.Lock()
		delete(processing, filePath)
		lock.Unlock()
	}
}
