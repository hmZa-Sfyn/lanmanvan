#!/usr/bin/env python3
"""
js-enum - Direct JS/HTML Secrets Scanner
- Finds ALL matches in a single URL
- Supports pattern=sk-* or pattern=./api_patterns.json
- Optional output=file.json to save results
- Prints count and all matches
"""

import os
import sys
import re
import json
import requests

USER_AGENT = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"

def load_patterns_from_file(filepath):
    try:
        with open(filepath, 'r') as f:
            raw = json.load(f)
        patterns = {}
        for name, regex in raw.items():
            try:
                patterns[name] = re.compile(regex)
            except re.error as e:
                print(f"[!] Invalid regex in '{name}': {e}", file=sys.stderr)
        return patterns
    except Exception as e:
        print(f"[!] Error loading pattern file '{filepath}': {e}", file=sys.stderr)
        sys.exit(1)

def get_builtin_patterns(pattern_str):
    if pattern_str == "sk-*":
        return {"sk-* secret": re.compile(r"sk-[a-zA-Z0-9]{20,}")}
    elif pattern_str == "eyJ*":
        return {"JWT-like": re.compile(r"eyJ[A-Za-z0-9_-]{10,}\\.[A-Za-z0-9_-]{10,}(?:\\.[A-Za-z0-9_-]{10,})?")}
    elif pattern_str == "AIza*":
        return {"Google API Key": re.compile(r"AIza[0-9A-Za-z\\-_]{35,}")}
    else:
        try:
            return {"Custom Regex": re.compile(pattern_str)}
        except re.error as e:
            print(f"[!] Invalid regex pattern: {e}", file=sys.stderr)
            sys.exit(1)

def scan_content(url, content, patterns):
    lines = content.splitlines()
    results = []

    for name, pat in patterns.items():
        for match in pat.finditer(content):
            secret = match.group(0)
            if len(secret) < 20:
                continue

            start = match.start()
            line_no = content[:start].count('\n') + 1
            line = lines[line_no - 1].strip() if 1 <= line_no <= len(lines) else ""

            result = {
                "service": name,
                "secret": secret,
                "url": url,
                "line_number": line_no,
                "line_snippet": line[:200]
            }
            results.append(result)

            # Print to stdout immediately
            print(f"⚠️  {name}")
            print(f"   → {secret}")
            print(f"   Line {line_no}: {line[:120]}{'...' if len(line) > 120 else ''}")

    return results

def main():
    url = os.getenv('ARG_URL')
    pattern_arg = os.getenv('ARG_PATTERN')
    output_file = os.getenv('ARG_OUTPUT', '').strip()

    if not url:
        print("[!] ARG_URL required", file=sys.stderr)
        sys.exit(1)
    if not pattern_arg:
        print("[!] ARG_PATTERN required (e.g., 'sk-*' or './api_patterns.json')", file=sys.stderr)
        sys.exit(1)

    url = url if url.startswith(('http://', 'https://')) else 'https://' + url

    # Load patterns
    if pattern_arg.endswith('.json') or pattern_arg == "./api_patterns.json":
        patterns = load_patterns_from_file(pattern_arg.replace('./', ''))
    else:
        patterns = get_builtin_patterns(pattern_arg)

    # Fetch content
    try:
        resp = requests.get(url, headers={'User-Agent': USER_AGENT}, timeout=15)
        resp.raise_for_status()
    except Exception as e:
        print(f"[!] Failed to fetch {url}: {e}", file=sys.stderr)
        sys.exit(1)

    print(f"Scanning: {url}")
    print(f"Pattern: {pattern_arg}")
    if output_file:
        print(f"Output: {output_file}")
    print("-" * 50)

    # Scan
    findings = scan_content(url, resp.text, patterns)

    # Summary
    print("-" * 50)
    print(f"[+] Total secrets found: {len(findings)}")

    # Save to file if requested
    if output_file and findings:
        try:
            # Ensure dir exists
            out_dir = os.path.dirname(output_file)
            if out_dir and not os.path.exists(out_dir):
                os.makedirs(out_dir)
            with open(output_file, 'w') as f:
                json.dump(findings, f, indent=2)
            print(f"[+] Results saved to: {output_file}")
        except Exception as e:
            print(f"[!] Failed to write output file: {e}", file=sys.stderr)
            sys.exit(1)
    elif output_file:
        print(f"[ ] No secrets found — output file not created.")

if __name__ == '__main__':
    main()