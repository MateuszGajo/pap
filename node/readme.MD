bytecode for js
node --print-bytecode --print-bytecode-filter=[funcName] filename.ext
assembly
node --allow-natives-syntax --print-opt-code --code-comments test.js


node profiling:
node --prof main.js ; node --prof-process isolate\*
Example:

```
 [Summary]:
   ticks  total  nonlib   name
      0    0.0%    0.0%  JavaScript
    428   64.8%  100.0%  C++
      5    0.8%    1.2%  GC
    233   35.2%          Shared libraries
```

Total can be more 100% if using multiple threads

perf record -F 99 node main.js ; perf report



emcc factorial.cpp -o factorial.wasm -s WASM=1 -s EXPORTED_FUNCTIONS="['_factorial','_QuadScalarPtr','_QuadAVXPtr','_process_array']" -s EXPORTED_RUNTIME_METHODS="['cwrap']" -msimd128  -msse3 -mavx2