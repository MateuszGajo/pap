#!/bin/bash

echo "Flushing filesystem buffers..."
sync

echo "Dropping caches..."
sudo sh -c 'echo 3 > /proc/sys/vm/drop_caches'

echo "Running Go program..."
./reader
