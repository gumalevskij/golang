#!/bin/bash

result_file="results.txt"
> "$result_file"

for parallel in 8 16 32 64 128 256 512; do
    make
    rm -rf database.json
    sed -i "s/parallel: [0-9]*/parallel: $parallel/" config.yaml
    result=$(./xkcd)
    echo "$parallel $result" >> "$result_file"
done

echo "Results in $result_file"