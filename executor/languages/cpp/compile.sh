#!/bin/sh
set -e
g++ main.cpp -O2 -static -s -o main 2> compile.stderr || true