package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/exp/mmap"
)

func readUsingNmap(filename string) int {
	reader, err := mmap.Open("./" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	// Get the file size
	fileSize := reader.Len()

	// Allocate a buffer and read the whole file
	p := make([]byte, reader.Len())

	_, err = reader.ReadAt(p, 0)

	if err != nil {
		panic(err)
	}

	return fileSize
}

func normalRead(filename string) int {
	file, err := os.ReadFile("./" + filename)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	return len(file)
}

func measure(count int, readerFunc func(filename string) int, name string, file string) {
	var start time.Time
	var size int
	var fileLoadTimeSec, fileSizeMb, readGBPerSec, minReadGbs, maxReadGbs float64

	for i := 0; i < count; i++ {

		start = time.Now()
		size = readerFunc(file)
		fileLoadTimeSec = time.Since(start).Seconds()
		fileSizeMb = float64(size / (1024 * 1024))
		readGBPerSec = fileSizeMb / (fileLoadTimeSec * 1000)
		if i == 0 {
			minReadGbs = readGBPerSec
			maxReadGbs = readGBPerSec
		} else if minReadGbs > readGBPerSec {
			minReadGbs = readGBPerSec
		} else if maxReadGbs < readGBPerSec {
			maxReadGbs = readGBPerSec
		}
	}

	fmt.Println("size", fileSizeMb)
	fmt.Println("time", fileLoadTimeSec)
	fmt.Printf("%v MIN: read file %v GB/s \n", name, minReadGbs)
	fmt.Printf("%v MAX: read file %v GB/s \n", name, maxReadGbs)
}

func main() {

	// measure(10, readUsingNmap, "NMAP", "midsize-output.json")
	// measure(10, normalRead, "BUILT-IN", "midsize-output.json")

	measure(4, readUsingNmap, "NMAP", "big-output.json")
	measure(4, normalRead, "BUILT-IN", "big-output.json")

}
