package cli

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// BuiltinFunction represents a builtin function with its handler
type BuiltinFunction struct {
	Name        string
	Description string
	Callback    func(args ...string) (string, error)
}

// BuiltinRegistry manages all builtin functions
type BuiltinRegistry struct {
	functions map[string]*BuiltinFunction
}

// NewBuiltinRegistry creates a new registry with all builtin functions
func NewBuiltinRegistry() *BuiltinRegistry {
	br := &BuiltinRegistry{
		functions: make(map[string]*BuiltinFunction),
	}
	br.registerAll()
	return br
}

// registerAll registers all builtin functions
func (br *BuiltinRegistry) registerAll() {
	// File system operations
	br.register("pwd", "Print current working directory", br.cmdPwd)
	br.register("cd", "Change directory (returns new path)", br.cmdCd)
	br.register("ls", "List directory contents", br.cmdLs)
	br.register("mkdir", "Create directory", br.cmdMkdir)
	br.register("rm", "Remove file or directory", br.cmdRm)
	br.register("cp", "Copy file or directory", br.cmdCp)
	br.register("mv", "Move or rename file", br.cmdMv)
	br.register("cat", "Read file contents", br.cmdCat)
	br.register("whoami", "Get current user", br.cmdWhoami)
	br.register("hostname", "Get system hostname", br.cmdHostname)
	br.register("date", "Get current date and time", br.cmdDate)
	br.register("uname", "Get system information", br.cmdUname)
	br.register("env", "Show environment variables", br.cmdEnv)
	br.register("which", "Find command location", br.cmdWhich)

	// Hash functions
	br.register("md5", "MD5 hash", br.cmdMd5)
	br.register("sha1", "SHA1 hash", br.cmdSha1)
	br.register("sha256", "SHA256 hash", br.cmdSha256)

	// Encoding/decoding
	br.register("base64", "Base64 encode/decode", br.cmdBase64)
	br.register("hex", "Hex encode/decode", br.cmdHex)
	br.register("url", "URL encode/decode", br.cmdUrl)
	br.register("json", "JSON format/minify", br.cmdJson)

	// String operations
	br.register("strlen", "Get string length", br.cmdStrlen)
	br.register("toupper", "Convert to uppercase", br.cmdToupper)
	br.register("tolower", "Convert to lowercase", br.cmdTolower)
	br.register("reverse", "Reverse string", br.cmdReverse)
	br.register("trim", "Trim whitespace", br.cmdTrim)

	// Network operations
	br.register("ping", "Ping host", br.cmdPing)
	br.register("nslookup", "DNS lookup", br.cmdNslookup)
	br.register("ipaddr", "Get IP addresses", br.cmdIpaddr)

	// Math operations
	br.register("calc", "Simple calculator", br.cmdCalc)

	// System commands
	br.register("sleep", "Sleep for seconds", br.cmdSleep)
	br.register("echo", "Print text", br.cmdEcho)
	br.register("readfile", "Read file and return content", br.cmdReadfile)

	// Utilities
	br.register("uuid", "Generate UUID", br.cmdUuid)
	br.register("timestamp", "Get current timestamp", br.cmdTimestamp)
	br.register("randomstr", "Generate random string", br.cmdRandomstr)
}

// register registers a single builtin function
func (br *BuiltinRegistry) register(name, desc string, callback func(args ...string) (string, error)) {
	br.functions[name] = &BuiltinFunction{
		Name:        name,
		Description: desc,
		Callback:    callback,
	}
}

// Execute runs a builtin function
func (br *BuiltinRegistry) Execute(name string, args ...string) (string, error) {
	fn, exists := br.functions[name]
	if !exists {
		return "", fmt.Errorf("builtin function '%s' not found", name)
	}
	return fn.Callback(args...)
}

// GetAll returns all registered functions
func (br *BuiltinRegistry) GetAll() map[string]*BuiltinFunction {
	return br.functions
}

// File system operations
func (br *BuiltinRegistry) cmdPwd(args ...string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return wd, nil
}

func (br *BuiltinRegistry) cmdCd(args ...string) (string, error) {
	if len(args) == 0 {
		home, _ := os.UserHomeDir()
		os.Chdir(home)
		return home, nil
	}
	if err := os.Chdir(args[0]); err != nil {
		return "", err
	}
	wd, _ := os.Getwd()
	return wd, nil
}

func (br *BuiltinRegistry) cmdLs(args ...string) (string, error) {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var output strings.Builder
	for _, entry := range entries {
		if entry.IsDir() {
			output.WriteString(fmt.Sprintf("%s/\n", entry.Name()))
		} else {
			output.WriteString(fmt.Sprintf("%s\n", entry.Name()))
		}
	}
	return strings.TrimSpace(output.String()), nil
}

func (br *BuiltinRegistry) cmdMkdir(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("mkdir requires a directory name")
	}
	if err := os.MkdirAll(args[0], 0755); err != nil {
		return "", err
	}
	return fmt.Sprintf("Directory '%s' created", args[0]), nil
}

func (br *BuiltinRegistry) cmdRm(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("rm requires a path")
	}
	if err := os.RemoveAll(args[0]); err != nil {
		return "", err
	}
	return fmt.Sprintf("Removed '%s'", args[0]), nil
}

func (br *BuiltinRegistry) cmdCp(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("cp requires source and destination")
	}
	source, dest := args[0], args[1]

	input, err := os.ReadFile(source)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(dest, input, 0644); err != nil {
		return "", err
	}
	return fmt.Sprintf("Copied '%s' to '%s'", source, dest), nil
}

func (br *BuiltinRegistry) cmdMv(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("mv requires source and destination")
	}
	source, dest := args[0], args[1]

	if err := os.Rename(source, dest); err != nil {
		return "", err
	}
	return fmt.Sprintf("Moved '%s' to '%s'", source, dest), nil
}

func (br *BuiltinRegistry) cmdCat(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("cat requires a file path")
	}

	content, err := os.ReadFile(args[0])
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (br *BuiltinRegistry) cmdWhoami(args ...string) (string, error) {
	user, err := os.LookupEnv("USER")
	if !err {
		user = "unknown"
	}
	return user, nil
}

func (br *BuiltinRegistry) cmdHostname(args ...string) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}

func (br *BuiltinRegistry) cmdDate(args ...string) (string, error) {
	format := "2006-01-02 15:04:05"
	if len(args) > 0 {
		format = args[0]
	}
	return time.Now().Format(format), nil
}

func (br *BuiltinRegistry) cmdUname(args ...string) (string, error) {
	cmd := exec.Command("uname", "-a")
	output, err := cmd.Output()
	if err != nil {
		return "unknown", nil
	}
	return strings.TrimSpace(string(output)), nil
}

func (br *BuiltinRegistry) cmdEnv(args ...string) (string, error) {
	var output strings.Builder
	for _, env := range os.Environ() {
		output.WriteString(env)
		output.WriteString("\n")
	}
	return strings.TrimSpace(output.String()), nil
}

func (br *BuiltinRegistry) cmdWhich(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("which requires a command name")
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return "", err
	}
	return path, nil
}

// Hash functions
func (br *BuiltinRegistry) cmdMd5(args ...string) (string, error) {
	return br.hash(md5.New(), args...)
}

func (br *BuiltinRegistry) cmdSha1(args ...string) (string, error) {
	return br.hash(sha1.New(), args...)
}

func (br *BuiltinRegistry) cmdSha256(args ...string) (string, error) {
	return br.hash(sha256.New(), args...)
}

func (br *BuiltinRegistry) hash(h hash.Hash, args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("hash requires input")
	}

	input := strings.Join(args, " ")
	io.WriteString(h, input)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Encoding/decoding
func (br *BuiltinRegistry) cmdBase64(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("base64 requires input")
	}

	input := strings.Join(args, " ")

	// Try to decode first
	if decoded, err := base64.StdEncoding.DecodeString(input); err == nil {
		// If it looks like valid base64, return decoded
		return string(decoded), nil
	}

	// Otherwise, encode
	return base64.StdEncoding.EncodeToString([]byte(input)), nil
}

func (br *BuiltinRegistry) cmdHex(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("hex requires input")
	}

	input := strings.Join(args, " ")

	// Try to decode first
	if decoded, err := hex.DecodeString(input); err == nil {
		return string(decoded), nil
	}

	// Otherwise, encode
	return hex.EncodeToString([]byte(input)), nil
}

func (br *BuiltinRegistry) cmdUrl(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("url requires input")
	}

	input := strings.Join(args, " ")

	// Simple URL encoding
	encoded := regexp.MustCompile(`[^a-zA-Z0-9\-_.]`).ReplaceAllStringFunc(input, func(s string) string {
		return fmt.Sprintf("%%%02X", s[0])
	})

	return encoded, nil
}

func (br *BuiltinRegistry) cmdJson(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("json requires input")
	}

	input := strings.Join(args, " ")
	var obj interface{}

	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return "", err
	}

	// Pretty print
	formatted, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

// String operations
func (br *BuiltinRegistry) cmdStrlen(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("strlen requires input")
	}

	input := strings.Join(args, " ")
	return strconv.Itoa(len(input)), nil
}

func (br *BuiltinRegistry) cmdToupper(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("toupper requires input")
	}

	input := strings.Join(args, " ")
	return strings.ToUpper(input), nil
}

func (br *BuiltinRegistry) cmdTolower(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("tolower requires input")
	}

	input := strings.Join(args, " ")
	return strings.ToLower(input), nil
}

func (br *BuiltinRegistry) cmdReverse(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("reverse requires input")
	}

	input := strings.Join(args, " ")
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

func (br *BuiltinRegistry) cmdTrim(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("trim requires input")
	}

	input := strings.Join(args, " ")
	return strings.TrimSpace(input), nil
}

// Network operations
func (br *BuiltinRegistry) cmdPing(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("ping requires a host")
	}

	host := args[0]
	timeout := time.Duration(5)
	if len(args) > 1 {
		if t, err := strconv.Atoi(args[1]); err == nil {
			timeout = time.Duration(t)
		}
	}

	// Try to resolve host and connect
	conn, err := net.DialTimeout("tcp", host+":80", timeout*time.Second)
	if err != nil {
		return fmt.Sprintf("Ping failed: %v", err), nil
	}
	defer conn.Close()

	return fmt.Sprintf("Ping to %s successful", host), nil
}

func (br *BuiltinRegistry) cmdNslookup(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("nslookup requires a hostname")
	}

	ips, err := net.LookupIP(args[0])
	if err != nil {
		return "", err
	}

	var output strings.Builder
	for _, ip := range ips {
		output.WriteString(ip.String())
		output.WriteString("\n")
	}

	return strings.TrimSpace(output.String()), nil
}

func (br *BuiltinRegistry) cmdIpaddr(args ...string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var output strings.Builder
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			output.WriteString(fmt.Sprintf("%s: %s\n", iface.Name, addr.String()))
		}
	}

	return strings.TrimSpace(output.String()), nil
}

// Math operations
func (br *BuiltinRegistry) cmdCalc(args ...string) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("calc requires: number operator number (e.g., '5 + 3')")
	}

	a, err1 := strconv.ParseFloat(args[0], 64)
	op := args[1]
	b, err2 := strconv.ParseFloat(args[2], 64)

	if err1 != nil || err2 != nil {
		return "", fmt.Errorf("invalid numbers")
	}

	var result float64
	switch op {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			return "", fmt.Errorf("division by zero")
		}
		result = a / b
	case "%":
		result = float64(int(a) % int(b))
	default:
		return "", fmt.Errorf("unknown operator: %s", op)
	}

	if result == float64(int(result)) {
		return strconv.Itoa(int(result)), nil
	}

	return strconv.FormatFloat(result, 'f', -1, 64), nil
}

// System commands
func (br *BuiltinRegistry) cmdSleep(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("sleep requires seconds")
	}

	seconds, err := strconv.Atoi(args[0])
	if err != nil {
		return "", err
	}

	time.Sleep(time.Duration(seconds) * time.Second)
	return fmt.Sprintf("Slept for %d seconds", seconds), nil
}

func (br *BuiltinRegistry) cmdEcho(args ...string) (string, error) {
	return strings.Join(args, " "), nil
}

func (br *BuiltinRegistry) cmdReadfile(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("readfile requires a file path")
	}

	content, err := os.ReadFile(args[0])
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// Utilities
func (br *BuiltinRegistry) cmdUuid(args ...string) (string, error) {
	// Simple UUID v4 generation using crypto/rand
	b := make([]byte, 16)
	for i := 0; i < len(b); i++ {
		// Simplified: using timestamp + counter for basic UUID-like string
		// In production, use crypto/rand
	}

	// Return a deterministic UUID for now
	t := time.Now().UnixNano()
	return fmt.Sprintf("550e8400-e29b-41d4-a716-%012d", t%1000000000000), nil
}

func (br *BuiltinRegistry) cmdTimestamp(args ...string) (string, error) {
	format := "unix"
	if len(args) > 0 {
		format = args[0]
	}

	t := time.Now()
	switch format {
	case "unix":
		return strconv.FormatInt(t.Unix(), 10), nil
	case "milli":
		return strconv.FormatInt(t.UnixMilli(), 10), nil
	case "nano":
		return strconv.FormatInt(t.UnixNano(), 10), nil
	default:
		return t.Format(format), nil
	}
}

func (br *BuiltinRegistry) cmdRandomstr(args ...string) (string, error) {
	length := 16
	if len(args) > 0 {
		if l, err := strconv.Atoi(args[0]); err == nil {
			length = l
		}
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var buffer bytes.Buffer

	for i := 0; i < length; i++ {
		// Use timestamp for seeding
		t := time.Now().UnixNano()
		idx := (t + int64(i)) % int64(len(charset))
		buffer.WriteByte(charset[idx])
	}

	return buffer.String(), nil
}
