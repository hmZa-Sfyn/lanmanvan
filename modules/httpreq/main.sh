#!/bin/bash
# HTTP Requests Module
# Makes HTTP requests and returns response

HOST="${ARG_HOST:-example.com}"
METHOD="${ARG_METHOD:-GET}"
PATH="${ARG_PATH:-/}"

echo "[*] Making $METHOD request to http://$HOST$PATH..."
echo

# Use curl with minimal output
if command -v curl &> /dev/null; then
    response=$(curl -s -X "$METHOD" "http://$HOST$PATH")
    if [ $? -eq 0 ]; then
        echo "[+] Response received:"
        echo "$response" | head -20
        if [ $(echo "$response" | wc -l) -gt 20 ]; then
            echo "[*] ... (output truncated)"
        fi
    else
        echo "[!] Request failed"
        exit 1
    fi
else
    echo "[!] curl not found"
    exit 1
fi

echo
