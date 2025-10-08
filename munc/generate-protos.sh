#!/bin/bash

set -e

cd proto

# Iterate through all files in the proto directory
for file in *; do
    # Check if the file has a .proto extension
    if [[ "$file" == *.proto ]]; then
        protoc --go_out=. "$file"
    fi
done

echo "Finished creating all proto files"