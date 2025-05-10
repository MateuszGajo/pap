function dualScalar(count, input) {
    let sumA = 0;
    let sumB = 0;
    for (let i = 0; i < count; i += 2) {
      sumA += input[i];
      sumB += input[i + 1];
    }
    const sum = sumA + sumB;
    return sum;
  }
  
  // Prepare the function for optimization
  %PrepareFunctionForOptimization(dualScalar);
  
  // Ensure the function will be optimized on the next call
  %OptimizeFunctionOnNextCall(dualScalar);
  
  // Call the function to trigger the optimization
  const input = [1, 2, 3, 4, 5, 6];
  const result = dualScalar(input.length, input);
  console.log(result);
  