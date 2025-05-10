rustc -C opt-level=3 main.rs  

-C opt-level=0: No optimizations. Fastest compile time, but slower execution.

-C opt-level=1: Basic optimizations that don’t take too long to apply.

-C opt-level=2 (default for release builds): More aggressive optimizations. It’s a good balance of compile time and performance.

-C opt-level=3: Maximum optimization, focuses on performance but can increase compile time.

-C opt-level=s: Optimizes for size (good for low-memory environments).

-C opt-level=z: Minimizes binary size even more than -C opt-level=s, but could reduce performance.