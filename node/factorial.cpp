#include <iostream>  
#include <immintrin.h> 
#include <cstdlib>
#include <string>
#include <climits>

typedef unsigned int u32;
u32 count = 4097;
u32* data[4096];

extern "C" {
    int factorial(int n) {
        if (n == 0) {
            return 1;
        } else {
            return n * factorial(n - 1);
        }
    }


    void __attribute__((target("avx2"))) preallocate(u32 Count, u32 *Input) {
        for (u32 i = 0; i < count; i++) {
            data[i] = &i;
        }
    }

    u32 __attribute__((target("avx2"))) QuadAVXPtr(u32 Count, u32 *Input)
{
    Count = count;
    Input = *data;
   
    u32* AlignedInput = (u32*)aligned_alloc(32, Count * sizeof(u32));

    if (AlignedInput == nullptr) {
        return 0;
    }

    for (u32 i = 0; i < Count; i++) {
        AlignedInput[i] = Input[i];
    }

    __m256i SumA = _mm256_setzero_si256();
    __m256i SumB = _mm256_setzero_si256();
    __m256i SumC = _mm256_setzero_si256();
    __m256i SumD = _mm256_setzero_si256();

    Count /= 32;
    while (Count--) {
        SumA = _mm256_add_epi32(SumA, _mm256_load_si256((__m256i *)&AlignedInput[0]));
        SumB = _mm256_add_epi32(SumB, _mm256_load_si256((__m256i *)&AlignedInput[8]));
        SumC = _mm256_add_epi32(SumC, _mm256_load_si256((__m256i *)&AlignedInput[16]));
        SumD = _mm256_add_epi32(SumD, _mm256_load_si256((__m256i *)&AlignedInput[24]));

        AlignedInput += 32;
    }

    __m256i SumAB = _mm256_add_epi32(SumA, SumB);
    __m256i SumCD = _mm256_add_epi32(SumC, SumD);
    __m256i Sum = _mm256_add_epi32(SumAB, SumCD);

    Sum = _mm256_hadd_epi32(Sum, Sum);
    Sum = _mm256_hadd_epi32(Sum, Sum);
    __m256i SumS = _mm256_permute2x128_si256(Sum, Sum, 1 | (1 << 4));
    Sum = _mm256_add_epi32(Sum, SumS);

    u32 result = _mm256_cvtsi256_si32(Sum);
    return result;
}  

    void process_array(int* arr, int size) {
        std::cout << "start process array" << std::endl;
        for (int i = 0; i < size; i++) {
            std::cout << "Array element " << i << ": " << arr[i] << std::endl;
        }
    }
    u32 QuadScalarPtr(u32 Count, u32 Input[4096])
    {
        u32 SumA = 0;
        u32 SumB = 0;
        u32 SumC = 0;
        u32 SumD = 0;
        
        count /= 4;
        while(count--)
        {
            SumA += data[0];
            SumB += data[1];
            SumC += data[2];
            SumD += data[3];
            data += 4;
        }
        u32 Sum = SumA + SumB + SumC + SumD;
        return Sum;
    }
}
int main() {
  int num = 5;
  std::cout << "Factorial of " << num << " is: " << factorial(num) << std::endl;
  return 0;
}
