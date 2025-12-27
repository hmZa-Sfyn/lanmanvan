#!/usr/bin/env python3
import os
import sys
import json
import requests

def main():
    url = os.getenv('ARG_URL', '').strip()
    key = os.getenv('ARG_KEY', '').strip()
    table = os.getenv('ARG_TABLE', '')
    
    if not url or not key:
        print("[!] Error: URL and key are required")
        sys.exit(1)
    
    # Ensure URL ends without trailing slash
    url = url.rstrip('/')
    
    headers = {
        'Authorization': f'Bearer {key}',
        'Content-Type': 'application/json',
        'apikey': key
    }
    
    try:
        print(f"[*] Testing Supabase connection to {url}")
        
        # Test basic connection
        response = requests.get(f'{url}/rest/v1/', headers=headers, timeout=10)
        
        if response.status_code in [200, 401, 403]:
            print(f"[+] Supabase endpoint is reachable!")
            print(f"[+] Status Code: {response.status_code}")
        else:
            print(f"[!] Unexpected status: {response.status_code}")
            sys.exit(1)
        
        # Try to list tables
        print(f"\n[*] Attempting to list tables...")
        tables_response = requests.get(f'{url}/rest/v1/information_schema.tables', 
                                      headers=headers, timeout=10)
        
        if tables_response.status_code == 200:
            print(f"[+] Successfully retrieved table information!")
            tables = tables_response.json()
            print(f"[+] Found {len(tables)} tables")
            for t in tables[:10]:
                print(f"    - {t.get('table_name', 'Unknown')}")
        else:
            print(f"[!] Could not retrieve tables: {tables_response.status_code}")
        
        # Test specific table if provided
        if table:
            print(f"\n[*] Testing access to table '{table}'...")
            table_response = requests.get(f'{url}/rest/v1/{table}?limit=1', 
                                         headers=headers, timeout=10)
            
            if table_response.status_code == 200:
                print(f"[+] Successfully accessed table '{table}'!")
                data = table_response.json()
                print(f"[+] Table has {len(data)} rows (showing first 1)")
            elif table_response.status_code == 404:
                print(f"[!] Table '{table}' not found")
            else:
                print(f"[!] Error accessing table: {table_response.status_code}")
        
        print(f"\n[+] Supabase connection test completed successfully!")
        
    except requests.exceptions.ConnectionError:
        print(f"[!] Error: Could not connect to Supabase URL")
        print(f"[!] Make sure the URL is correct: {url}")
        sys.exit(1)
    except Exception as e:
        print(f"[!] Error: {e}")
        sys.exit(1)

if __name__ == '__main__':
    main()
