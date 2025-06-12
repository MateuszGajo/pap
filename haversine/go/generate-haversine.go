package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
)

type DataStruct struct {
	Pairs []HarversineData `json:"pairs"`
}

type HarversineData struct {
	X0 float64 `json:"x0"`
	X1 float64 `json:"x1"`
	Y0 float64 `json:"y0"`
	Y1 float64 `json:"y1"`
}

func seedData() {
	dataCount := 50_000_000

	haversineData := make([]HarversineData, dataCount)

	for i := 0; i < dataCount; i++ {
		haversineData[i] = HarversineData{X0: rand.Float64()*360 - 180, X1: rand.Float64()*360 - 180, Y0: rand.Float64()*180 - 90, Y1: rand.Float64()*180 - 90}
	}

	data := DataStruct{
		Pairs: haversineData,
	}

	jsonOutput, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("output.json")

	if err != nil {
		log.Fatal(err)
	}

	file.Write(jsonOutput)
}
