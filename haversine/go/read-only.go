package main

import (
	"log"
	"os"
)

func read() []byte {

	data, err := os.ReadFile("./output.json")

	if err != nil {
		log.Fatal(err)
	}
	return data
}
