#!/usr/bin/env python3
import os
import sys
import json
import requests

def main():
    url = os.getenv('ARG_URL', '').strip()
    data_str = os.getenv('ARG_DATA', '').strip()
    headers_str = os.getenv('ARG_HEADERS', '')
    
    if not url or not data_str:
        print("[!] Error: URL and data are required")
        sys.exit(1)
    
    headers = {'Content-Type': 'application/json'}
    if headers_str:
        try:
            headers.update(json.loads(headers_str))
        except:
            print("[!] Error: Invalid JSON headers")
            sys.exit(1)
    
    # Try to parse data as JSON, otherwise send as string
    try:
        data = json.loads(data_str)
    except:
        data = data_str
    
    try:
        print(f"[*] Sending POST request to {url}")
        response = requests.post(url, json=data if isinstance(data, dict) else None, 
                               data=data if isinstance(data, str) else None,
                               headers=headers, timeout=10)
        
        print(f"[+] Status Code: {response.status_code}")
        print(f"[+] Response Length: {len(response.content)} bytes")
        print(f"\n[*] Response:")
        print(response.text[:2000])
        if len(response.text) > 2000:
            print(f"    ... (truncated, total {len(response.text)} chars)")
        
    except Exception as e:
        print(f"[!] Error: {e}")
        sys.exit(1)

if __name__ == '__main__':
    main()
