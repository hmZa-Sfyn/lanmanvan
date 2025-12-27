#!/usr/bin/env python3
import os
import sys
import json
import requests

def main():
    url = os.getenv('ARG_URL', '').strip()
    headers_str = os.getenv('ARG_HEADERS', '')
    timeout = int(os.getenv('ARG_TIMEOUT', '10'))
    
    if not url:
        print("[!] Error: URL is required")
        sys.exit(1)
    
    headers = {}
    if headers_str:
        try:
            headers = json.loads(headers_str)
        except:
            print("[!] Error: Invalid JSON headers")
            sys.exit(1)
    
    try:
        print(f"[*] Sending GET request to {url}")
        response = requests.get(url, headers=headers, timeout=timeout, allow_redirects=True)
        
        print(f"[+] Status Code: {response.status_code}")
        print(f"[+] Content Length: {len(response.content)} bytes")
        print(f"[+] Content Type: {response.headers.get('Content-Type', 'Unknown')}")
        print(f"\n[*] Response Headers:")
        for key, value in response.headers.items():
            print(f"    {key}: {value}")
        
        print(f"\n[*] Response Body:")
        print(response.text[:2000])
        if len(response.text) > 2000:
            print(f"    ... (truncated, total {len(response.text)} chars)")
        
    except requests.exceptions.Timeout:
        print(f"[!] Error: Request timed out after {timeout}s")
        sys.exit(1)
    except requests.exceptions.ConnectionError as e:
        print(f"[!] Error: Connection failed - {e}")
        sys.exit(1)
    except Exception as e:
        print(f"[!] Error: {e}")
        sys.exit(1)

if __name__ == '__main__':
    main()
