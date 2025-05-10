use std::arch::x86_64::_rdtsc;

type U32 = u32;

fn single_scalar(count: U32, input: &[U32]) -> U32 {
    let mut sum: U32 = 0;
    for i in 0..count {
        sum += input[i as usize];
    }
    sum
}

fn unroll2_scalar(count: U32, input: &[U32]) -> U32 {
    let mut sum: U32 = 0;
    let mut i = 0;
    while i + 1 < count {
        sum += input[i as usize];
        sum += input[(i + 1) as usize];
        i += 2;
    }
    sum
}

fn unroll4_scalar(count: U32, input: &[U32]) -> U32 {
    let mut sum: U32 = 0;
    let mut i = 0;
    while i + 3 < count {
        sum += input[i as usize];
        sum += input[(i + 1) as usize];
        sum += input[(i + 2) as usize];
        sum += input[(i + 3) as usize];
        i += 4;
    }
    sum
}

fn dual_scalar(count: U32, input: &[U32]) -> U32 {
    let mut sum_a: U32 = 0;
    let mut sum_b: U32 = 0;
    let mut i = 0;
    while i + 1 < count {
        sum_a += input[i as usize];
        sum_b += input[(i + 1) as usize];
        i += 4;
    }
    let sum: u32 = sum_a + sum_b;
    sum
}

fn quad_scalar(count: U32, input: &[U32]) -> U32 {
    let mut sum_a: U32 = 0;
    let mut sum_b: U32 = 0;
    let mut sum_c: U32 = 0;
    let mut sum_d: U32 = 0;
    let mut i = 0;
    while i + 3 < count {
        sum_a += input[i as usize];
        sum_b += input[(i + 1) as usize];
        sum_c += input[(i + 2) as usize];
        sum_d += input[(i + 3) as usize];
        i += 4;
    }
    let sum: u32 = sum_a + sum_b+ sum_c + sum_d;
    sum
}


fn benchmark<F>(label: &str, func: F, data: &[U32], count: U32, test_count: usize, clock_ghz: f64)
where
    F: Fn(U32, &[U32]) -> U32,
{
    let mut min_cycles = u64::MAX;
    let mut res = 0;

    for _ in 0..test_count {
        let start = unsafe { _rdtsc() };
        res = func(count, data);
        let end = unsafe { _rdtsc() };
        let cycles = end - start;

        if cycles < min_cycles {
            min_cycles = cycles;
        }
    }

    let seconds = min_cycles as f64 / (clock_ghz * 1e9);
    let cycles_per_add = min_cycles as f64 / count as f64;
    let adds_per_cycle = count as f64 / min_cycles as f64;

    println!("=== {} ===", label);
    println!("Sum: {}", res);
    println!("CPU Cycles Taken: {}", min_cycles);
    println!("Time (seconds): {}", seconds);
    println!("cycle/adds: {}", cycles_per_add);
    println!("adds/cycle: {}", adds_per_cycle);
}

fn main() {
    let clock_ghz = 2.9; // GHz
    let count: U32 = 4097;
    let test_count = 10_000;
    let mut data = vec![0; count as usize];
    for i in 0..count {
        data[i as usize] = i;
    }

    benchmark("SingleScalar", single_scalar, &data, count, test_count, clock_ghz);
    benchmark("Unroll2Scalar", unroll2_scalar, &data, count, test_count, clock_ghz);
    benchmark("Unroll4Scalar", unroll4_scalar, &data, count, test_count, clock_ghz);
    benchmark("dualScalar", dual_scalar, &data, count, test_count, clock_ghz);
    benchmark("quadScalar", quad_scalar, &data, count, test_count, clock_ghz);
}
