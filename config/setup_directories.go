package config

import (
	"log"
	"os"
)

func SetupDirectories() {
	if _, err := os.Stat(App().ErrorParsePath); os.IsNotExist(err) {
		err := os.Mkdir(App().ErrorParsePath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(App().PathToWatch); os.IsNotExist(err) {
		err := os.Mkdir(App().PathToWatch, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}
