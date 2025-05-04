import ctypes
import os

here = os.path.dirname(__file__)
lib_path = os.path.join(here, "libtsc.so")
tsc_lib = ctypes.CDLL(lib_path)
tsc_lib.read_tsc.restype = ctypes.c_uint64


def single_scalar(count, input):
    sum = 0
    index = 0
    while index + 1 < count:
        sum += input[index]
        index += 1
    return sum

def unroll2_scalar(count, input):
    sum = 0
    index = 0
    while index + 1 < count:
        sum += input[index] + input[index + 1]
        index += 2
    return sum

def unroll4_scalar(count, input):
    sum = 0
    index = 0
    while index + 1 < count:
        sum += input[index] + input[index + 1] + input[index + 2]  + input[index + 3] 
        index += 4
    return sum

def dual_scalar(count, input):
    sumA = 0
    sumB = 0
    index = 0
    while index + 1 < count:
        sumA += input[index]
        sumB +=  input[index + 1]
        index += 2
    sum = sumA + sumB
    return sum

def quad_scalar(count, input):
    sumA = 0
    sumB = 0
    sumC = 0
    sumD = 0
    index = 0
    while index + 1 < count:
        sumA += input[index]
        sumB +=  input[index + 1]
        sumC +=  input[index + 2]
        sumD +=  input[index + 3]
        index += 4
    sum = sumA + sumB + sumC + sumD
    return sum

def benchmark(name, func, count, data, clock_ghz, test_count= 10000):
    best_cycles = float('inf')
    result = 0


    for _ in range(test_count):
        start = tsc_lib.read_tsc()
        result = func(count, data)
        end = tsc_lib.read_tsc()
        cycles = end - start
        if 0 < cycles < best_cycles :
            best_cycles = cycles

    seconds = best_cycles / (clock_ghz * 1e9)
    cycles_per_add = best_cycles / count
    adds_per_cycle = count / best_cycles

    print(f"=== {name} ===")
    print(f"Sum: {result}")
    print(f"CPU Cycles Taken: {best_cycles}")
    print(f"Time (seconds): {seconds}")
    print(f"Cycle/Adds: {cycles_per_add}")
    print(f"Adds/Cycle: {adds_per_cycle}\n")

def main():
    clock_ghz = 4.2
    count = 4096
    data = list(range(1, count + 1))

    benchmark("SingleScalar", single_scalar, count, data, clock_ghz)
    benchmark("Unroll2Scalar", unroll2_scalar, count, data, clock_ghz)
    benchmark("Unroll4Scalar", unroll4_scalar, count, data, clock_ghz)
    benchmark("dualScalar", dual_scalar, count, data, clock_ghz)
    benchmark("quadScalar", quad_scalar, count, data, clock_ghz)

if __name__ == "__main__":
    main()