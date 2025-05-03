// tsc.c
#include <stdint.h>
#include <stdio.h>

// Inline assembly to call rdtsc on x86 processors
uint64_t read_tsc() {
    uint64_t tsc;
    __asm__ volatile ("rdtsc" : "=A"(tsc));
    return tsc;
}
