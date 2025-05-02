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
import "fmt"


func singleScalar(count int, input []int) int {
	sum := 0

	for index := 0; index < count; index++ {
		sum += input[index]
	}

	return sum
}

func main() {
	clockGHz := 4.2
	var count int = 4096
	data := make([]int, count)

	for i := 1; i < count; i++ {
		data[i] = i
	}

	start := C.read_tsc()
	res := singleScalar(count, data)

	end :=  C.read_tsc()
	cycles := int64(end - start)
	// cycles := int64(1)
	seconds := float64(cycles) / (clockGHz * 1e9)
	cyclesPerAdd := float64(cycles) / float64(count)
	addsPerCycle := float64(count) / float64(cycles)

	fmt.Println("Sum: " ,res)
	fmt.Println("CPU Cycles Taken: " ,cycles)
	fmt.Println("Time (seconds): " ,seconds)
	fmt.Println("cycle/adds: " ,cyclesPerAdd)
	fmt.Println("adds/cycle: " ,addsPerCycle)

}