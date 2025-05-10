package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

var fileLoadTime float64
var calculateTime float64

func loadHaversine() DataStruct {
	start := time.Now()
	fileData, err := os.ReadFile("./output.json")

	if err != nil {
		log.Fatal(err)
	}

	var data DataStruct

	err = json.Unmarshal(fileData, &data)

	if err != nil {
		log.Fatal(err)
	}
	fileLoadTime = time.Since(start).Seconds()

	return data
}

var radian = float64(math.Pi / 180)

func haversineOfDegrees(x0, y0, x1, y1, r float64) float64 {
	dY := (y1 - y0) * radian
	dX := (x1 - x0) * radian
	y0 = (y0) * radian
	y1 = (y1) * radian

	RootTerm := math.Pow(math.Sin(dY/2), 2) + math.Cos(y0)*math.Cos(y1)*math.Pow(math.Sin(dX/2), 2)
	Result := 2 * r * math.Asin(math.Sqrt(RootTerm))

	return Result
}

func calculateHaversine(data DataStruct) (float64, int) {
	start := time.Now()

	earthRadiuskm := 6371.0
	sum := 0.0
	count := 0

	for _, v := range data.Pairs {
		sum += haversineOfDegrees(v.X0, v.Y0, v.X1, v.Y1, earthRadiuskm)
		count += 1
	}

	average := sum / float64(count)
	calculateTime = time.Since(start).Seconds()

	return average, count
}

func harverstineWork() {
	data := loadHaversine()
	avg, count := calculateHaversine(data)

	fmt.Println("Result: ", avg)
	fmt.Println("Input = ", fileLoadTime, " seconds")
	fmt.Println("Math = ", calculateTime, " seconds")
	fmt.Println("Total = ", calculateTime+fileLoadTime, " seconds")
	fmt.Println("count", count)
	fmt.Println("Throughput = ", float64(count)/(calculateTime+fileLoadTime), " haversines/second")
}
