package winsys

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
)

var (
	privNames     = make(map[string]int64)
	privNameMutex sync.Mutex
)

func SetTokenPrivilege(i *Inject) error {
	localTokenHandle, err := getToken(i)
	if err != nil {
		return err
	}

	//First pass to get token privs
	_, err = GetTokenPrivileges(localTokenHandle)
	if err != nil {
		return err
	}

	// Start LookupPrivilegeValue
	// https://docs.microsoft.com/en-us/windows/desktop/secauthz/privilege-constants
	var tokenPrivilege []string
	if len(tokenPrivilege) == 0 {
		tokenPrivilege = append(tokenPrivilege, SE_DEBUG_NAME)
	} else {
		tokenPrivilege = strings.Split(strings.Replace(i.Privilege, " ", "", -1), ",")
	}

	mpvn, err := MapPrivilegeValueToName(tokenPrivilege)
	if err != nil {
		return errors.Wrap(err, "Error: Could not get LUID from privilege")
	}
	//byte representation of the luid and attributes struct
	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, uint32(len(mpvn)))
	var luid int64
	for _, p := range mpvn {
		luid = p
		binary.Write(&b, binary.LittleEndian, p)
		binary.Write(&b, binary.LittleEndian, uint32(SE_PRIVILEGE_ENABLED))
	}

	success, err := AdjustTokenPrivileges(localTokenHandle, false, &b.Bytes()[0], uint32(b.Len()), nil, nil)
	if success == 0 {
		return err
	}
	if err == ERROR_NOT_ALL_ASSIGNED {
		return err
	}
	fmt.Printf("Successfully added %s to LUID %v\n", tokenPrivilege[0], luid)

	//Second pass for getting token privs (see if they stuck)
	_, err = GetTokenPrivileges(localTokenHandle)
	if err != nil {
		return err
	}
	return nil
}

func getToken(i *Inject) (syscall.Token, error) {
	currentProcessHandle, err := syscall.GetCurrentProcess()
	if err != nil {
		return 0, err
	}

	// fmt.Printf("Current Process Handle: %v\n", currentProcessHandle)
	err = syscall.OpenProcessToken(currentProcessHandle, TOKEN_ADJUST_PRIVILEGES|TOKEN_QUERY, &i.Token.tokenHandle)
	if err != nil {
		return 0, err
	}
	return i.Token.tokenHandle, nil
}

func LookupPrivilegeValue(systemName string, name string, luid *int64) (err error) {
	var sN *uint16
	sN, err = syscall.UTF16PtrFromString(systemName)
	if err != nil {
		return
	}
	var n *uint16
	n, err = syscall.UTF16PtrFromString(name)
	if err != nil {
		return
	}
	r1, _, e1 := syscall.Syscall(ProcLookupPrivilegeValueW.Addr(), 3, uintptr(unsafe.Pointer(sN)), uintptr(unsafe.Pointer(n)), uintptr(unsafe.Pointer(luid)))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func MapPrivilegeValueToName(names []string) ([]int64, error) {
	var privileges []int64
	privNameMutex.Lock()
	defer privNameMutex.Unlock()
	for _, name := range names {
		p, ok := privNames[name]
		if !ok {
			err := LookupPrivilegeValue("", name, &p)
			if err != nil {
				return nil, errors.Wrapf(err, "LookupPrivilegeValue failed on '%v'", name)
			}
			privNames[name] = p
		}
		privileges = append(privileges, p)
	}
	return privileges, nil
}

func LookupPrivilegeName(systemName string, luid int64) (string, error) {
	buf := make([]uint16, 256)
	bufSize := uint32(len(buf))
	var sN *uint16
	sN, err := syscall.UTF16PtrFromString(systemName)
	if err != nil {
		return "", err
	}
	syscall.StringToUTF16Ptr(systemName)
	r1, _, e1 := syscall.Syscall6(ProcLookupPrivilegeNameW.Addr(), 4, uintptr(unsafe.Pointer(sN)), uintptr(unsafe.Pointer(&luid)), uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&bufSize)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return syscall.UTF16ToString(buf), nil
}

func GetTokenPrivileges(token syscall.Token) (map[string]Privilege, error) {
	var size uint32
	syscall.GetTokenInformation(token, syscall.TokenPrivileges, nil, 0, &size)
	b := bytes.NewBuffer(make([]byte, size))
	err := syscall.GetTokenInformation(token, syscall.TokenPrivileges, &b.Bytes()[0], uint32(b.Len()), &size)
	if err != nil {
		return nil, err
	}
	var privilegeCount uint32
	err = binary.Read(b, binary.LittleEndian, &privilegeCount)
	if err != nil {
		return nil, err
	}
	fmt.Printf("[+] Identified %v privileges...\n", privilegeCount)
	rtn := make(map[string]Privilege, privilegeCount)
	for i := 0; i < int(privilegeCount); i++ {
		var luid int64
		err = binary.Read(b, binary.LittleEndian, &luid)
		if err != nil {
			return nil, err
		}

		var attributes uint32
		err = binary.Read(b, binary.LittleEndian, &attributes)
		if err != nil {
			return nil, err
		}

		name, err := LookupPrivilegeName("", luid)
		if err != nil {
			return nil, err
		}
		// https://docs.microsoft.com/en-us/windows/desktop/api/winnt/ns-winnt-_privilege_set
		rtn[name] = Privilege{
			LUID:             luid,
			Name:             name,
			EnabledByDefault: (attributes & SE_PRIVILEGE_ENABLED_BY_DEFAULT) > 0,
			Enabled:          (attributes & SE_PRIVILEGE_ENABLED) > 0,
			Removed:          (attributes & SE_PRIVILEGE_REMOVED) > 0,
			Used:             (attributes & SE_PRIVILEGE_USED_FOR_ACCESS) > 0,
		}
	}
	for k, v := range rtn {
		fmt.Println("Priv Name (key): ", k)
		fmt.Println("Priv Enabled: ", v.Enabled)
		fmt.Println("Priv Enabled by Default: ", v.EnabledByDefault)
		fmt.Println("LUID: ", v.LUID)
		fmt.Println("Priv Name: ", v.Name)
		fmt.Println("Priv Removed: ", v.Removed)
		fmt.Println("Priv Used: ", v.Used)
		fmt.Printf("\n\n")
	}
	return rtn, nil
}

// GetTokenUser returns the User associated with the given Token.
func GetTokenUser(token syscall.Token) (User, error) {
	tokenUser, err := token.GetTokenUser()
	if err != nil {
		return User{}, err
	}

	var user User
	user.SID, err = tokenUser.User.Sid.String()
	if err != nil {
		return user, err
	}

	user.Account, user.Domain, user.Type, err = tokenUser.User.Sid.LookupAccount("")
	if err != nil {
		return user, err
	}

	fmt.Println("Account: ", user.Account)
	fmt.Println("Domain: ", user.Domain)
	// https://docs.microsoft.com/en-us/windows/desktop/cimwin32prov/win32-account
	fmt.Println("SID Type: ", user.Type)
	fmt.Println("SID: ", user.SID)

	return user, nil
}

// adjustTokenPrivileges from core zsyscall_windows.go
func AdjustTokenPrivileges(token syscall.Token, disableAllPrivileges bool, newstate *byte, buflen uint32, prevstate *byte, returnlen *uint32) (ret uint32, err error) {
	var _p0 uint32
	if disableAllPrivileges {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r0, _, e1 := syscall.Syscall6(ProcAdjustTokenPrivileges.Addr(), 6, uintptr(token), uintptr(_p0), uintptr(unsafe.Pointer(newstate)), uintptr(buflen), uintptr(unsafe.Pointer(prevstate)), uintptr(unsafe.Pointer(returnlen)))
	ret = uint32(r0)
	if true {
		if e1 != 0 {
			err = e1
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
