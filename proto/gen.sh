#!/bin/bash

echo "Starting generate proto"

# get dependency
./get-deps.sh

# get googleapi lib
./get-googleapi.sh

# compile proto
./get-compile.sh

echo "Successfully generated proto"