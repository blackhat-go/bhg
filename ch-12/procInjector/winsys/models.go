package winsys

import "syscall"

type Inject struct {
	Pid              uint32
	DllPath          string
	DLLSize          uint32
	Privilege        string
	RemoteProcHandle uintptr
	Lpaddr           uintptr
	LoadLibAddr      uintptr
	RThread          uintptr
	Token            TOKEN
}

type Privilege struct {
	LUID             int64
	Name             string
	EnabledByDefault bool
	Enabled          bool
	Removed          bool
	Used             bool
}

type TOKEN struct {
	tokenHandle syscall.Token
}

// User represent the information about a Windows account.
type User struct {
	SID     string
	Account string
	Domain  string
	Type    uint32
}
