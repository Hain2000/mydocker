package container

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := os.Pipe() // cmd在readPipe读取数据
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init") // /proc/self/exe 调用自身初始化环境
	// fork 新进程时，通过指定 Cloneflags 会创建对应的 Namespace 以实现隔离，这里包括UTS（主机名）、PID（进程ID）、挂载点、网络、IPC等方面的隔离。
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe} //  让cmd使用readPipe FD
	cmd.Dir = "/root/busybox"
	return cmd, writePipe
}
