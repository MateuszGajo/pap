#include <string>
#include <iostream>
#include <x86intrin.h>
#include <climits>
#include <chrono>
#include <future>
#include <numeric>     
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

u32 __attribute__((target("ssse3"))) SingleSSE(u32 Count, u32 *Input) {
	__m128i Sum = _mm_setzero_si128();
	for(u32 Index = 0; Index < Count; Index += 4)
	{
		Sum = _mm_add_epi32(Sum, _mm_load_si128((__m128i *)&Input[Index]));
	}

	Sum = _mm_hadd_epi32(Sum, Sum);
	Sum = _mm_hadd_epi32(Sum, Sum);
	
	return _mm_cvtsi128_si32(Sum);
}

u32 __attribute__((target("avx2"))) SingleAVX(u32 Count, u32 *Input)
{
	__m256i Sum = _mm256_setzero_si256();
	for(u32 Index = 0; Index < Count; Index += 8)
	{
		Sum = _mm256_add_epi32(Sum, _mm256_loadu_si256((__m256i *)&Input[Index]));
	}

	Sum = _mm256_hadd_epi32(Sum, Sum);
	Sum = _mm256_hadd_epi32(Sum, Sum);
	__m256i SumS = _mm256_permute2x128_si256(Sum, Sum, 1 | (1 << 4));
	Sum = _mm256_add_epi32(Sum, SumS);
	
	return _mm256_cvtsi256_si32(Sum);
}

u32 __attribute__((target("avx2"))) DualAVX(u32 Count, u32 *Input)
{
	__m256i SumA = _mm256_setzero_si256();
	__m256i SumB = _mm256_setzero_si256();
	for(u32 Index = 0; Index < Count; Index += 16)
	{
		SumA = _mm256_add_epi32(SumA, _mm256_loadu_si256((__m256i *)&Input[Index]));
		SumB = _mm256_add_epi32(SumB, _mm256_loadu_si256((__m256i *)&Input[Index + 8]));
	}

	__m256i Sum = _mm256_add_epi32(SumA, SumB);

	Sum = _mm256_hadd_epi32(Sum, Sum);
	Sum = _mm256_hadd_epi32(Sum, Sum);
	__m256i SumS = _mm256_permute2x128_si256(Sum, Sum, 1 | (1 << 4));
	Sum = _mm256_add_epi32(Sum, SumS);
	
	return _mm256_cvtsi256_si32(Sum);
}

u32 __attribute__((target("avx2"))) QuadAVX(u32 Count, u32 *Input)
{
	__m256i SumA = _mm256_setzero_si256();
	__m256i SumB = _mm256_setzero_si256();
	__m256i SumC = _mm256_setzero_si256();
	__m256i SumD = _mm256_setzero_si256();
	for(u32 Index = 0; Index < Count; Index += 32)
	{
		SumA = _mm256_add_epi32(SumA, _mm256_loadu_si256((__m256i *)&Input[Index]));
		SumB = _mm256_add_epi32(SumB, _mm256_loadu_si256((__m256i *)&Input[Index + 8]));
		SumC = _mm256_add_epi32(SumC, _mm256_loadu_si256((__m256i *)&Input[Index + 16]));
		SumD = _mm256_add_epi32(SumD, _mm256_loadu_si256((__m256i *)&Input[Index + 24]));
	}

	__m256i SumAB = _mm256_add_epi32(SumA, SumB);
	__m256i SumCD = _mm256_add_epi32(SumC, SumD);
	__m256i Sum = _mm256_add_epi32(SumAB, SumCD);

	Sum = _mm256_hadd_epi32(Sum, Sum);
	Sum = _mm256_hadd_epi32(Sum, Sum);
	__m256i SumS = _mm256_permute2x128_si256(Sum, Sum, 1 | (1 << 4));
	Sum = _mm256_add_epi32(Sum, SumS);
	
	return _mm256_cvtsi256_si32(Sum);
}

u32 __attribute__((target("avx2"))) QuadAVXPtr(u32 Count, u32 *Input)
{
	__m256i SumA = _mm256_setzero_si256();
	__m256i SumB = _mm256_setzero_si256();
	__m256i SumC = _mm256_setzero_si256();
	__m256i SumD = _mm256_setzero_si256();
	
	Count /= 32;
	while(Count--)
	{
		SumA = _mm256_add_epi32(SumA, _mm256_loadu_si256((__m256i *)&Input[0]));
		SumB = _mm256_add_epi32(SumB, _mm256_loadu_si256((__m256i *)&Input[8]));
		SumC = _mm256_add_epi32(SumC, _mm256_loadu_si256((__m256i *)&Input[16]));
		SumD = _mm256_add_epi32(SumD, _mm256_loadu_si256((__m256i *)&Input[24]));
		
		Input += 32;
	}

	__m256i SumAB = _mm256_add_epi32(SumA, SumB);
	__m256i SumCD = _mm256_add_epi32(SumC, SumD);
	__m256i Sum = _mm256_add_epi32(SumAB, SumCD);

	Sum = _mm256_hadd_epi32(Sum, Sum);
	Sum = _mm256_hadd_epi32(Sum, Sum);
	__m256i SumS = _mm256_permute2x128_si256(Sum, Sum, 1 | (1 << 4));
	Sum = _mm256_add_epi32(Sum, SumS);
	
	return _mm256_cvtsi256_si32(Sum);
}


u32 SingleScalarForThrea(u32 start, u32 end, u32* Input) {
	u32 Sum = 0;

	for (u32 Index = start; Index < end; ++Index) {
		Sum += Input[Index];
	}

	return Sum;
}

u32 SingleScalarMultiThread(u32 Count, u32 *Input)
{

	std::future<u32> ret1 = std::async(&SingleScalarForThrea, (u32)0,Count / 4, Input);
	std::future<u32> ret2 = std::async(&SingleScalarForThrea, Count / 4,2*Count / 4, Input);
	std::future<u32> ret3 = std::async(&SingleScalarForThrea, 2*Count / 4,3*Count / 4, Input);
	std::future<u32> ret4 = std::async(&SingleScalarForThrea, 3*Count / 4,Count, Input);

	u32 total = ret1.get() + ret2.get() + ret3.get() + ret4.get();

	return total;
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
	benchmark("SingleSSE", SingleSSE, data, count, testCount, clockGHz);
	benchmark("SingleAVX", SingleAVX, data, count, testCount, clockGHz);
	benchmark("DualAVX", DualAVX, data, count, testCount, clockGHz);
	benchmark("QuadAVX", QuadAVX, data, count, testCount, clockGHz);
	benchmark("QuadAVXPtr", QuadAVXPtr, data, count, testCount, clockGHz);

	const u32 biggerDataCountl2 = 32768;
	u32* biggerDatal2 = new u32[biggerDataCountl2];

	for (u32 i = 0; i < biggerDataCountl2; i++) {
		biggerDatal2[i] = i;
	}

	benchmark("QuadAVXPtr(32768 size)", QuadAVXPtr, biggerDatal2, biggerDataCountl2, testCount, clockGHz);

	const u32 biggerDataCountl3 = 2611444;
	u32* biggerDatal3 = new u32[biggerDataCountl3];
	for (u32 i = 0; i < biggerDataCountl3; i++) {
		biggerDatal3[i] = i;
	}
	benchmark("SingleScalar", SingleScalar,biggerDatal3, biggerDataCountl3, testCount, clockGHz);
	benchmark("SingleScalarMultiThrea", SingleScalarMultiThread, biggerDatal3, biggerDataCountl3, testCount, clockGHz);
	benchmark("QuadAVXPtr(2611444 size)", QuadAVXPtr, biggerDatal3, biggerDataCountl3, testCount, clockGHz);

	// const u32 biggerDataCountOutsidel3 = 33554432;
	// u32* biggerDataOutsidel3 = new u32[biggerDataCountOutsidel3];
	// for (u32 i = 0; i < biggerDataCountOutsidel3; i++) {
	// 	biggerDataOutsidel3[i] = i;
	// }

	// benchmark("QuadAVXPtr(33554432 size)", QuadAVXPtr, biggerDataOutsidel3, biggerDataCountOutsidel3, testCount, clockGHz);


	return 0;
}
