const { performance } = require("perf_hooks");
const fs = require('fs');
const path = require('path');
const { WASI } = require('wasi');

const clockGHz = 2.9;

function singleScalar(count, input) {
  let sum = 0;
  for (let i = 0; i < count; i++) {
    sum += input[i];
  }
  return sum;
}

function unroll2Scalar(count, input) {
  let sum = 0;
  for (let i = 0; i < count; i += 2) {
    sum += input[i];
    sum += input[i + 1];
  }
  return sum;
}

function unroll4Scalar(count, input) {
  let sum = 0;
  for (let i = 0; i < count; i += 4) {
    sum += input[i];
    sum += input[i + 1];
    sum += input[i + 2];
    sum += input[i + 3];
  }
  return sum;
}

function dualScalar(count, input) {
  let sumA = 0;
  let sumB = 0;
  for (let i = 0; i < count; i += 2) {
    sumA += input[i];
    sumB += input[i + 1];
  }
  const sum = sumA + sumB
  return sum;
}

function benchmark(name, func, data, count, testCount) {

  let elapsedTimeMs = 100000000;
  let res;
  for (i = 0; i < testCount; i++) {
    const start = performance.now();
    res = func(count, data);
    const end = performance.now();
    const currentElapsedTimeMs = end - start;
    if (elapsedTimeMs > currentElapsedTimeMs) {
      elapsedTimeMs = currentElapsedTimeMs
    }
  }

  // Treat it as ballpark figure, there is no equivalent for rdtsc in node
  const cycles = elapsedTimeMs * clockGHz * 1e6;

  const cyclesPerAdd = cycles / count;
  const addsPerCycle = count / cycles;

  console.log(`======== ${name} ======`);
  console.log(`Sum: ${res}`);
  console.log(`Elapsed Time (ms): ${elapsedTimeMs}`);
  console.log(`Estimated CPU Cycles: ${cycles}`);
  console.log(`Cycles per Addition: ${cyclesPerAdd}`);
  console.log(`Additions per Cycle: ${addsPerCycle}`);

}



async function main() {
  const wasi = new WASI({
    version: 'preview1', // This is required in Node.js 20+
    args: process.argv,
    env: process.env
  });

  const wasmBuffer = fs.readFileSync(path.resolve('factorial.wasm'));
  const module = await WebAssembly.compile(wasmBuffer);
  const instance = await WebAssembly.instantiate(module, {
    wasi_snapshot_preview1: wasi.wasiImport
  });

  wasi.start(instance);
  const funcs = instance.exports;
  const count = 4096;
  const testCount = 10000;
  let data = Array.from({ length: count }, (_, i) => i + 1);

  benchmark("singleScalar", singleScalar, data, count, testCount)
  benchmark("unroll2Scalar", unroll2Scalar, data, count, testCount)
  benchmark("unroll4Scalar", unroll4Scalar, data, count, testCount)
  benchmark("dualScalar", dualScalar, data, count, testCount)


  data = new Uint32Array(count);

  for (let i = 0; i < count; i++) {
    data[i] = i + 1;
  }
  console.log("===== u32 array =====")
  benchmark("singleScalar", singleScalar, data, count, testCount)
  benchmark("unroll2Scalar", unroll2Scalar, data, count, testCount)
  benchmark("unroll4Scalar", unroll4Scalar, data, count, testCount)
  benchmark("dualScalar", dualScalar, data, count, testCount)
  // Degradation of performance using uint32
  // For further improvements lets use webassembly https://nodejs.org/en/learn/getting-started/nodejs-with-webassembly

  // One way to do it, todo: improve it later
  funcs.preallocate()
  benchmark("QuadScalarPtr via webassembly 0-4096", funcs.QuadScalarPtr, data, count, testCount)
  benchmark("QuadAVXPtr via webassembly 0-4096", funcs.QuadAVXPtr, data, count, testCount)

}


main();
