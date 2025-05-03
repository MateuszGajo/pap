import ctypes
import os

lib_path = os.path.abspath('./tsc.dll')
tsc_lib = ctypes.CDLL(lib_path)
tsc_lib.read_tsc.restype = ctypes.c_uint64


def single_scalar(count, input):
    sum = 0
    for index in range(count):
        sum += input[index]
    return sum


def main():
    clock_ghz = 4.2
    count = 4096
    data = list(range(1, count + 1)) 

    start = tsc_lib.read_tsc()

    res = single_scalar(count, data)

    end = tsc_lib.read_tsc()
    cycles = end - start
    seconds = float(cycles) / (clock_ghz * 1e9)
    cycles_per_add = float(cycles) / float(count)
    adds_per_cycle = float(count) / float(cycles)

    print(f"Sum: {res}")
    print(f"CPU Cycles Taken: {cycles}")
    print(f"Time (seconds): {seconds}")
    print(f"Cycle/Adds: {cycles_per_add}")
    print(f"Adds/Cycle: {adds_per_cycle}")

if __name__ == "__main__":
    main()