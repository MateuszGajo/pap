/* ========================================================================
   LISTING 1
   ======================================================================== */
#include <string>
#include <iostream>
#include <intrin.h>

using namespace std;


typedef unsigned int u32;
u32 SingleScalar(u32 Count, u32* Input) {
	u32 Sum = 0;

	for (u32 Index = 0; Index < Count; ++Index) {
		Sum +=Input[Index];
	}

	return Sum;
}



int __declspec(noinline) add(int A, int B)
{
	return A + B;
}

#pragma optimize("", off)
int main(int ArgCount, char** Args)
{
	unsigned long long start, end;
	const double clockGHz = 4.2;// assume turbo speed  
	const int count = 4096;    
	u32 data[count];

	for (u32 i = 1; i < count; i++) {
		data[i] = i;
	}

	
	start = __rdtsc();
	u32 res = SingleScalar(count, data);

	end = __rdtsc();
	unsigned __int64 cycles = end - start;
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