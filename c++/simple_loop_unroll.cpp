#include <string>
#include <iostream>
#include <x86intrin.h>
#include <climits>
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

	for (u32 Index = 0; Index < Count; Index += 4) {
		Sum += Input[Index];
		Sum += Input[Index + 1];
		Sum += Input[Index + 2];
		Sum += Input[Index + 3];
	}

	return Sum;
}

u32 DualScalar(u32 Count, u32* Input) {
	u32 SumA = 0;
	u32 SumB = 0;


	for (u32 Index=0; Index < Count; Index +=2) {
		SumA += Input[Index];
		SumB += Input[Index+1];
	}
	u32 Sum = SumA + SumB;
	return Sum;
}

u32 QuadScalar(u32 Count, u32* Input) {
	u32 SumA = 0;
	u32 SumB = 0;
	u32 SumC = 0;
	u32 SumD = 0;


	for (u32 Index=0; Index < Count; Index +=4) {
		SumA += Input[Index];
		SumB += Input[Index+1];
		SumC += Input[Index+2];
		SumD += Input[Index+3];
	}
	u32 Sum = SumA + SumB + SumC + SumD;
	return Sum;
}

u32 QuadScalarPtr(u32 Count, u32 *Input)
{
	u32 SumA = 0;
	u32 SumB = 0;
	u32 SumC = 0;
	u32 SumD = 0;
	
	Count /= 4;
	while(Count--)
	{
		SumA += Input[0];
		SumB += Input[1];
		SumC += Input[2];
		SumD += Input[3];
		Input += 4;
	}
	
	u32 Sum = SumA + SumB + SumC + SumD;
	return Sum;
}

void benchmark(const string& label, u32(*func)(u32,u32*),u32* data, u32 count, int testCount, double clockGHz) {
	unsigned long long start, end;

	u_int64_t cycles = INT_MAX;
	u32 res = 0;

	for (int i = 0; i < testCount; i++) {
		start = __rdtsc();
		res = func(count, data);

		end = __rdtsc();
		u_int64_t  currentCycle = end - start;
		if (currentCycle < cycles) {
			cycles = currentCycle;
		};
	}
	double seconds = cycles / (clockGHz * 1e9);
	double cyclesPerAdd = static_cast<double>(cycles) / count;
	double addsPerCycle = static_cast<double>(count) / cycles;

	cout << "=== " << label << " ===" << '\n';
	cout << "Sum: " << res << '\n';
	cout << "CPU Cycles Taken: " << cycles << '\n';
	cout << "Time (seconds): " << seconds << '\n';
	cout << "cycle/adds " << cyclesPerAdd << '\n';
	cout << "adds/cycle " << addsPerCycle << '\n';

}

int main(int ArgCount, char** Args)
{
	const double clockGHz = 2.9;  // Base speed, turbo mode should be disable for predictable results
	const u32 count = 4097;
	const int testCount = 10000;
	u32 data[count];

	for (u32 i = 0; i < count; i++) {
		data[i] = i;
	}

	benchmark("SingleScalar", SingleScalar, data, count, testCount, clockGHz);
	benchmark("Unroll2Scalar", Unroll2Scalar, data, count, testCount, clockGHz);
	benchmark("Unroll4Scalar", Unroll4Scalar, data, count, testCount, clockGHz);
	benchmark("DualScalar", DualScalar, data, count, testCount, clockGHz);
	benchmark("QuadScalar", QuadScalar, data, count, testCount, clockGHz);
	benchmark("QuadScalarPtr", QuadScalarPtr, data, count, testCount, clockGHz);


	return 0;
}
