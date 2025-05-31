package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"time"

	"github.com/valyala/fastjson"
)

var fileLoadTimeSec float64
var calculateTime float64

func loadHaversineWithFastJson() DataStruct {
	start := time.Now()
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	var data DataStruct
	file, err := os.ReadFile("./output.json")
	var p fastjson.Parser
	_, err = p.ParseBytes(file)
	if err != nil {
		panic(err)
	}
	pprof.StopCPUProfile()
	fileLoadTimeSec = time.Since(start).Seconds()
	fileSizeMb := float64(len(file) / (1024 * 1024))
	fmt.Println("file load time plus parsing ", fileLoadTimeSec)
	fmt.Println("file size mb: ", fileSizeMb)
	fmt.Printf("read file + parse json %v mb/s \n", fileSizeMb/fileLoadTimeSec)
	fmt.Printf("read file + parse json %v GB/s \n", fileSizeMb/(fileLoadTimeSec*1024))
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func loadHaversineWithBuildInMethod() DataStruct {
	start := time.Now()
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	var data DataStruct
	file, err := os.ReadFile("./output.json")
	json.Unmarshal(file, &data)
	pprof.StopCPUProfile()
	fileLoadTimeSec = time.Since(start).Seconds()
	fileSizeMb := float64(len(file) / (1024 * 1024))
	fmt.Println("file load time plus parsing ", fileLoadTimeSec)
	fmt.Println("file size mb: ", fileSizeMb)
	fmt.Printf("read file + parse json %v mb/s \n", fileSizeMb/fileLoadTimeSec)
	fmt.Printf("read file + parse json %v GB/s \n", fileSizeMb/(fileLoadTimeSec*1024))
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func loadHaversine() DataStruct {
	start := time.Now()
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}

	fileData, size := parseJson()

	pprof.StopCPUProfile()

	var data DataStruct
	fileLoadTimeSec = time.Since(start).Seconds()
	fileSizeMb := float64(size / (1024 * 1024))
	fmt.Println("file load time plus parsing ", fileLoadTimeSec)
	fmt.Println("file size mb: ", fileSizeMb)
	fmt.Printf("read file + parse json %v mb/s \n", fileSizeMb/fileLoadTimeSec)
	fmt.Printf("read file + parse json %v GB/s \n", fileSizeMb/(fileLoadTimeSec*1024))
	err = assign(fileData, &data)

	if err != nil {
		log.Fatal(err)
	}

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
	fmt.Println("Input = ", fileLoadTimeSec, " seconds")
	fmt.Println("Math = ", calculateTime, " seconds")
	fmt.Println("Total = ", calculateTime+fileLoadTimeSec, " seconds")
	fmt.Println("count", count)
	fmt.Println("Throughput = ", float64(count)/(calculateTime+fileLoadTimeSec), " haversines/second")
}

func harverstineWorkBuiltIn() {
	data := loadHaversineWithBuildInMethod()
	avg, count := calculateHaversine(data)

	fmt.Println("Result: ", avg)
	fmt.Println("Input = ", fileLoadTimeSec, " seconds")
	fmt.Println("Math = ", calculateTime, " seconds")
	fmt.Println("Total = ", calculateTime+fileLoadTimeSec, " seconds")
	fmt.Println("count", count)
	fmt.Println("Throughput = ", float64(count)/(calculateTime+fileLoadTimeSec), " haversines/second")
}

func harverstineWorkFastJson() {
	data := loadHaversineWithFastJson()
	avg, count := calculateHaversine(data)

	fmt.Println("Result: ", avg)
	fmt.Println("Input = ", fileLoadTimeSec, " seconds")
	fmt.Println("Math = ", calculateTime, " seconds")
	fmt.Println("Total = ", calculateTime+fileLoadTimeSec, " seconds")
	fmt.Println("count", count)
	fmt.Println("Throughput = ", float64(count)/(calculateTime+fileLoadTimeSec), " haversines/second")
}
