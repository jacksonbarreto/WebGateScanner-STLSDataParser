package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/config"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/processor"
	"github.com/jacksonbarreto/WebGateScanner-STLSDataParser/internal/sslresponseparser"
	"github.com/jacksonbarreto/WebGateScanner-kafka/producer"
	"log"
	"strings"
	"sync"
)

func main() {
	const defaultBufferSize = 100
	const configFilePath = ""
	config.InitConfig(configFilePath)
	config.SetupDirectories()

	kafkaProducer, producerErr := producer.NewProducer(config.Kafka().TopicProducer,
		config.Kafka().Brokers, config.Kafka().MaxRetry)
	if producerErr != nil {
		panic(producerErr)
	}
	defer func(kafkaProducer *producer.Producer) {
		err := kafkaProducer.Close()
		if err != nil {
			// TODO: log error
		}
	}(kafkaProducer)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			// TODO: log error
		}
	}(watcher)

	processFileQueueSize := config.App().ProcessFileQueueSize
	if processFileQueueSize == 0 {
		processFileQueueSize = defaultBufferSize
	}
	filesToProcess := make(chan string, processFileQueueSize)
	totalWorkers := config.App().Workers
	var filesInProcess = make(map[string]bool)
	var lock = sync.Mutex{}

	for i := 0; i < totalWorkers; i++ {
		fileProcessor := processor.NewDefaultFileProcessor(kafkaProducer,
			sslresponseparser.NewDefaultAssessmentSSLResponseParser(),
			processor.NewDefaultHostExtractor(config.App().ProcessFileExtension),
			config.App().ReadyToProcessSuffix, &lock, filesInProcess)
		go fileProcessor.ProcessFileFromChannel(filesToProcess)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create && strings.HasSuffix(event.Name, config.App().ReadyToProcessSuffix) {
					log.Println("NewProcessor file detected:", event.Name)
					originalFileName := strings.TrimSuffix(event.Name, config.App().ReadyToProcessSuffix)
					lock.Lock()
					if !filesInProcess[event.Name] {
						filesInProcess[event.Name] = false
						filesToProcess <- originalFileName
					}
					lock.Unlock()
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(config.App().PathToWatch)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
