#!/bin/bash

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root (for dropping caches)." 
   exit 1
fi

echo "=== CLEANUP: Dropping Linux caches (cold start) ==="
sync; echo 3 > /proc/sys/vm/drop_caches
echo "Caches dropped."

# First run to warm up
/usr/bin/time -v ./program 
echo "second time"
# Second run with actual measurement
/usr/bin/time -v ./program