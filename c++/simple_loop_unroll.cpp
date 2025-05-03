#include <string>
#include <iostream>
#include <intrin.h>
using namespace std;


typedef unsigned int u32;
u32 SingleScalar(u32 Count, u32* Input) {
	u32 Sum = 0;

	for (u32 Index = 0; Index < Count; ++Index) {
		Sum += Input[Index];
	}

	return Sum;
}

u32 Unroll2Scalar(u32 Count, u32* Input) {
	u32 Sum = 0;

	for (u32 Index = 0; Index < Count; Index += 2) {
		Sum += Input[Index];
		Sum += Input[Index + 1];
	}

	return Sum;
}

u32 Unroll4Scalar(u32 Count, u32* Input) {
	u32 Sum = 0;

	for (u32 Index = 0; Index < Count; Index + 4) {
		Sum += Input[Index];
		Sum += Input[Index + 1];
		Sum += Input[Index + 2];
		Sum += Input[Index + 3];
	}

	return Sum;
}



//#pragma optimize("", off)
int main(int ArgCount, char** Args)
{
	unsigned long long start, end;
	const double clockGHz = 2.9;// assume turbo speed  
	const int count = 4097;
	const int testCount = 10000;
	u32 data[count];

	for (u32 i = 0; i < count; i++) {
		data[i] = i;
	}

	unsigned __int64 cycles = INT_MAX;
	u32 res = 0;

	for (int i = 0; i < testCount; i++) {
		start = __rdtsc();
		//res = SingleScalar(count, data);
		res = Unroll2Scalar(count, data);

		end = __rdtsc();
		unsigned __int64 currentCycle = end - start;
		if (currentCycle < cycles) {
			cycles = currentCycle;
		};
	}

	double seconds = cycles / (clockGHz * 1e9);
	double cyclesPerAdd = static_cast<double>(cycles) / count;
	double addsPerCycle = static_cast<double>(count) / cycles;

	cout << "Sum: " << res << '\n';
	cout << "CPU Cycles Taken: " << cycles << '\n';
	cout << "Time (seconds): " << seconds << '\n';
	cout << "cycle/adds " << cyclesPerAdd << '\n';
	cout << "adds/cycle " << addsPerCycle << '\n';

	return 0;
}
