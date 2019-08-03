package winsys

import "syscall"

var (
	ModKernel32 = syscall.NewLazyDLL("kernel32.dll")
	modUser32   = syscall.NewLazyDLL("user32.dll")
	modAdvapi32 = syscall.NewLazyDLL("Advapi32.dll")

	ProcOpenProcessToken      = modAdvapi32.NewProc("GetProcessToken")
	ProcLookupPrivilegeValueW = modAdvapi32.NewProc("LookupPrivilegeValueW")
	ProcLookupPrivilegeNameW  = modAdvapi32.NewProc("LookupPrivilegeNameW")
	ProcAdjustTokenPrivileges = modAdvapi32.NewProc("AdjustTokenPrivileges")
	ProcGetAsyncKeyState      = modUser32.NewProc("GetAsyncKeyState")
	ProcVirtualAlloc          = ModKernel32.NewProc("VirtualAlloc")
	ProcCreateThread          = ModKernel32.NewProc("CreateThread")
	ProcWaitForSingleObject   = ModKernel32.NewProc("WaitForSingleObject")
	ProcVirtualAllocEx        = ModKernel32.NewProc("VirtualAllocEx")
	ProcVirtualFreeEx         = ModKernel32.NewProc("VirtualFreeEx")
	ProcCreateRemoteThread    = ModKernel32.NewProc("CreateRemoteThread")
	ProcGetLastError          = ModKernel32.NewProc("GetLastError")
	ProcWriteProcessMemory    = ModKernel32.NewProc("WriteProcessMemory")
	ProcOpenProcess           = ModKernel32.NewProc("OpenProcess")
	ProcGetCurrentProcess     = ModKernel32.NewProc("GetCurrentProcess")
	ProcIsDebuggerPresent     = ModKernel32.NewProc("IsDebuggerPresent")
	ProcGetProcAddress        = ModKernel32.NewProc("GetProcAddress")
	ProcCloseHandle           = ModKernel32.NewProc("CloseHandle")
	ProcGetExitCodeThread     = ModKernel32.NewProc("GetExitCodeThread")
)
