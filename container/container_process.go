package container

import (
	"os"
	"os/exec"
	"syscall"
)

func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}              // 这里的 init 指令就用用来在子进程中调用 initCommand
	cmd := exec.Command("/proc/self/exe", args...) // /proc/self/exe 调用自身初始化环境

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
	return cmd
}
