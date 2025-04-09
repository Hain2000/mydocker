package container

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
)

func RunContainerInitProcess(command string, args []string) error {
	log.Infof("command:%s", command)
	_ = unix.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// MS_NOEXEC 在本文件系统 许运行其程序。
	// MS_NOSUID 在本系统中运行程序的时候， 允许 set-user-ID set-group-ID
	// MS_NOD 所有 mount 的系统都会默认设定的参数。

	// 相当于 mount -t proc proc /proc
	// -t proc 指定挂载的文件系统类型是 proc
	// 第二个 proc	挂载源（source），是一个名字，也可以写成 none
	// 挂在到 /proc目录
	_ = unix.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}
