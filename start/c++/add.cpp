///* ========================================================================
//   LISTING 1
//   ======================================================================== */
//#include <string>
//#include <iostream>
//#include <intrin.h>
//#include <chrono>
//#include <thread>
//#include <cstdint>
//using namespace std;
//
//typedef unsigned int u32;
//u32 SingleScalar(u32 Count, u32* Input) {
//	u32 Sum = 0;
//
//	for (u32 Index = 0; Index < Count; ++Index) {
//		Sum += Input[Index];
//	}
//
//	return Sum;
//}
//
//
//
//int __declspec(noinline) add(int A, int B)
//{
//	return A + B;
//}
//
//double MeasureCPUFrequencyGHz(int milliseconds = 100) {
//	using namespace std::chrono;
//
//	auto start_time = high_resolution_clock::now();
//	uint64_t start_cycles = __rdtsc();
//
//	std::this_thread::sleep_for(std::chrono::milliseconds(milliseconds));
//	uint64_t end_cycles = __rdtsc();
//	auto end_time = high_resolution_clock::now();
//
//	duration<double> elapsed_seconds = end_time - start_time;
//	double elapsed_time = elapsed_seconds.count(); // in seconds
//
//	uint64_t cycles = end_cycles - start_cycles;
//	double GHz = static_cast<double>(cycles) / (elapsed_time * 1e9);
//
//	return GHz;
//}
//
//
//#pragma optimize("", off)
//int main(int ArgCount, char** Args)
//{
//	unsigned long long start, end;
//	const double clockGHz = 4.2;// assume turbo speed  
//	const int count = 4096;
//	u32 data[count];
//	const int testCount = 10000;
//	 
//	for (u32 i = 1; i < count; i++) {
//		data[i] = i;
//	}
//
//
//	unsigned __int64 cycles = INT_MAX;
//	u32 res = 0;
//	double ghz = MeasureCPUFrequencyGHz();
//	cout << "Measured CPU Clock: " << ghz << " GHz" << endl;
//
//	//	start = __rdtsc();
//	//for (int i = 0; i < testCount; i++) {
//	//	//u32 res = SingleScalar(count, data);
//	//	//res = Unroll2Scalar(count, data);
//
//	//	end = __rdtsc();
//	//	uint64_t currentCycle = end - start;
//	//	if (currentCycle < cycles) {
//	//		cycles = currentCycle;
//	//	};
//	//}
//	double seconds = cycles / (clockGHz * 1e9);
//	double cyclesPerAdd = static_cast<double>(cycles) / count;
//	double addsPerCycle = static_cast<double>(count) / cycles;
//
//	cout << "Sum: " << res << '\n';
//	cout << "CPU Cycles Taken: " << cycles << '\n';
//	cout << "Time (seconds): " << seconds << '\n';
//	cout << "cycle/adds " << cyclesPerAdd << '\n';
//	cout << "adds/cycle " << addsPerCycle << '\n';
//
//	return 0;
//}