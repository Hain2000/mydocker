package main

import (
	log "github.com/sirupsen/logrus"
	"mydocker/cgroups"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
	"os"
	"strings"
)

// Run 执行具体 command
/*
这里的Start方法是真正开始执行由NewParentProcess构建好的command的调用，它首先会clone出来一个namespace隔离的
进程，然后在子进程中，调用/proc/self/exe,也就是调用自己，发送init参数，调用我们写的init方法，
去初始化容器的一些资源。
*/
func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, volume, containerName, imageName string, envSlice []string) {
	containerId := container.GenerateContainerID() // 生成 10 位容器 id
	parent, writePipe := container.NewParentProcess(tty, volume, containerId, imageName, envSlice)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Errorf("Run parent.Start err:%v", err)
	}

	// record container info
	err := container.RecordContainerInfo(parent.Process.Pid, comArray, containerName, containerId, volume)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}

	// 创建cgroup manager, 并通过调用set和apply设置资源限制并使限制在容器上生效
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	_ = cgroupManager.Set(res)
	_ = cgroupManager.Apply(parent.Process.Pid, res)

	sendInitCommand(comArray, writePipe)
	if tty { // 如果是tty，那么父进程等待，就是前台运行，否则就是跳过，实现后台运行
		_ = parent.Wait()
		container.DeleteWorkSpace(containerId, volume)
		container.DeleteContainerInfo(containerId)
	}
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ") // 把命令拼起来
	log.Infof("command all is %s", command)
	_, _ = writePipe.WriteString(command) // 然后写到管道里
	_ = writePipe.Close()
}
