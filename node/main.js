const { performance } = require("perf_hooks");

const clockGHz = 4.2;

function singleScalar(count, input) {
  let sum = 0;
  for (let i = 0; i < count; i++) {
    sum += input[i];
  }
  return sum;
}

function main() {
  const count = 4096;
  const data = Array.from({ length: count }, (_, i) => i + 1);
  const start = performance.now();
  const res = singleScalar(count, data);
  const end = performance.now();

  const elapsedTimeMs = end - start;

  // Treat it as ballpark figure, there is no equivalent for rdtsc in node
  const cycles = elapsedTimeMs * clockGHz * 1e6;

  const cyclesPerAdd = cycles / count;
  const addsPerCycle = count / cycles;

  console.log(`Sum: ${res}`);
  console.log(`Elapsed Time (ms): ${elapsedTimeMs}`);
  console.log(`Estimated CPU Cycles: ${cycles}`);
  console.log(`Cycles per Addition: ${cyclesPerAdd}`);
  console.log(`Additions per Cycle: ${addsPerCycle}`);
}

main();
