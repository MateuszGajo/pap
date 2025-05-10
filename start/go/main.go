package main

/*
#include <stdint.h>

// Inline assembly to call __rdtsc on x86 processors
static inline uint64_t read_tsc() {
    uint64_t tsc;
    // Inline assembly for reading TSC (Time Stamp Counter)
    __asm__ volatile ("rdtsc" : "=A"(tsc));
    return tsc;
}
*/
import "C"
import (
	"fmt"
	"math"
)

func singleScalar(count int, input []int) int {
	sum := 0

	for index := 0; index < count; index++ {
		sum += input[index]
	}

	return sum
}

func unroll2Scalar(count int, input []int) int {
	sum := 0

	for index := 0; index < count; index += 2 {
		sum += input[index]
		sum += input[index+1]
	}

	return sum
}

func unroll4Scalar(count int, input []int) int {
	sum := 0

	for index := 0; index < count; index += 4 {
		sum += input[index]
		sum += input[index+1]
		sum += input[index+2]
		sum += input[index+3]
	}

	return sum
}

func dualScalar(count int, input []int) int {
	sumA := 0
	sumB := 0

	for index := 0; index < count; index += 2 {
		sumA += input[index]
		sumB += input[index+1]
	}

	sum := sumA + sumB
	return sum
}

func quadScalar(count int, input []int) int {
	sumA := 0
	sumB := 0
	sumC := 0
	sumD := 0

	for index := 0; index < count; index += 4 {
		sumA += input[index]
		sumB += input[index+1]
		sumC += input[index+2]
		sumD += input[index+3]
	}

	sum := sumA + sumB + sumC + sumD
	return sum
}

func benchmark(name string, data []int, function func(count int, input []int) int, count int, testCount int, clockGHz float64) {

	var res int
	var start, end float64
	var cycles = math.MaxFloat64

	for i := 0; i < testCount; i++ {
		start = float64(C.read_tsc())

		res = function(count, data)
		end = float64(C.read_tsc())
		currentCycle := float64(end - start)
		if currentCycle < cycles {
			cycles = currentCycle
		}
	}

	seconds := float64(cycles) / (clockGHz * 1e9)
	cyclesPerAdd := float64(cycles) / float64(count)
	addsPerCycle := float64(count) / float64(cycles)

	fmt.Println("====== ", name, " ========")
	fmt.Println("Sum: ", res)
	fmt.Println("CPU Cycles Taken: ", cycles)
	fmt.Println("Time (seconds): ", seconds)
	fmt.Println("cycle/adds: ", cyclesPerAdd)
	fmt.Println("adds/cycle: ", addsPerCycle)
}

func main() {
	clockGHz := 4.2
	count := 4096
	testCount := 10000
	data := make([]int, count)

	for i := 1; i < count; i++ {
		data[i] = i
	}

	benchmark("singleScalar", data, singleScalar, count, testCount, clockGHz)
	benchmark("unroll2Scalar", data, unroll2Scalar, count, testCount, clockGHz)
	benchmark("unroll4Scalar", data, unroll4Scalar, count, testCount, clockGHz)
	benchmark("dualScalar", data, dualScalar, count, testCount, clockGHz)
	benchmark("quadScalar", data, quadScalar, count, testCount, clockGHz)
}
