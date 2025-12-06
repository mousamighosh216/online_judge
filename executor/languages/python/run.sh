#!/bin/sh
ulimit -t 2
ulimit -v 65536
if [ -x /box/usr/bin/python3 ]; then
    /box/usr/bin/python3 main.py
else 
    # fallback
    /usr/bin/python3 main.py
fi