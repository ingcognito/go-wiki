#!/bin/sh
echo "Building GO Binary"
go build -o main . && ./main