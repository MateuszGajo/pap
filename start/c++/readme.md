// compile no optimization
g++ -O0 -g -o my_program simple_loop_unroll.cpp

SIMD - simple instruction multiple data
SAFE A LOT OF FRONT END WORK - meaning we dont need to figure out that multiple add are dependant on each other and can't be execute in parallel, so having instruction PADDD we know that we can execute few adds at once
SSE - simi streaming extenstion 4 bytes
AUX - 8 bytes
AUX-512 -16 bytes
PADDD - pack add (d word)