package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
	"os"
)

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit container to image",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing image name")
		}
		containerID := context.Args().Get(0)
		imageName := context.Args().Get(1)
		return commitContainer(containerID, imageName)
	},
}

var listCommand = cli.Command{
	Name:  "ps",
	Usage: "list all the containers",
	Action: func(context *cli.Context) error {
		ListContainers()
		return nil
	},
}

var logCommand = cli.Command{
	Name:  "logs",
	Usage: "print logs of a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("please input your container name")
		}
		containerName := context.Args().Get(0)
		logContainer(containerName)
		return nil
	},
}

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
			mydocker run -it [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it", // 简单起见，这里把 -i 和 -t 参数合并成一个
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "mem", // 限制进程内存使用量，为了避免和 stress 命令的 -m 参数冲突 这里使用 -mem,到时候可以看下解决冲突的方法
			Usage: "memory limit,e.g.: -mem 100m",
		},
		cli.StringFlag{
			Name:  "memswap", // 限制进程swap，系统在内存压力下 优先使用 swap 的程度。该值越高，系统在内存紧张时越倾向于将内存页交换到 swap 中，反之则越倾向于保留内存中的数据。
			Usage: "memory limit,e.g.: -memswap 0",
		},
		cli.StringFlag{
			Name:  "cpu",
			Usage: "cpu quota,e.g.: -cpu 100", // 限制进程 cpu 使用率
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit,e.g.: -cpuset 2,4", // 限制进程 cpu 使用率
		},
		cli.StringFlag{ // 数据卷
			Name:  "v",
			Usage: "volume,e.g.: -v /ect/conf:/etc/conf",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name，e.g.: -name mycontainer",
		},
		cli.StringSliceFlag{ // 增加 -e flag
			Name:  "e",
			Usage: "set environment,e.g. -e name=mydocker",
		},
	},
	/*
		这里是run命令执行的真正函数。
		1.判断参数是否包含command
		2.获取用户指定的command
		3.调用Run function去准备启动容器:
	*/
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		imageName := cmdArray[0]
		cmdArray = cmdArray[1:]

		tty := context.Bool("it")
		detach := context.Bool("d")
		if tty && detach {
			return fmt.Errorf("it and d paramter can not both provided")
		}
		if !detach { // 如果不是指定后台运行，就默认前台运行
			tty = true
		}
		resConf := &subsystems.ResourceConfig{
			MemoryLimit:     context.String("mem"),
			MemorySwapLimit: context.Int("memswap"),
			CpuSet:          context.String("cpuset"),
			CpuCfsQuota:     context.Int("cpu"),
		}
		log.Info("resConf:", resConf)
		volume := context.String("v")
		containerName := context.String("name")
		envSlice := context.StringSlice("e")
		Run(tty, cmdArray, resConf, volume, containerName, imageName, envSlice)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		err := container.RunContainerInitProcess()
		return err
	},
}

var execCommand = cli.Command{
	Name:  "exec",
	Usage: "exec a command into container",
	Action: func(context *cli.Context) error {
		// 如果环境变量存在，说明C代码已经运行过了，即setns系统调用已经执行了，这里就直接返回，避免重复执行
		if os.Getenv(EnvExecPid) != "" {
			log.Infof("pid callback pid %v", os.Getgid())
			return nil
		}
		// 格式：mydocker exec 容器名字 命令，因此至少会有两个参数
		if len(context.Args()) < 2 {
			return fmt.Errorf("missing container name or command")
		}
		containerName := context.Args().Get(0)
		// 将除了容器名之外的参数作为命令部分
		commandArray := context.Args().Tail()
		ExecContainer(containerName, commandArray)
		return nil
	},
}

var stopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop a container,e.g. mydocker stop 1234567890",
	Action: func(context *cli.Context) error {
		// 期望输入是：mydocker stop 容器Id，如果没有指定参数直接打印错误
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container id")
		}
		containerId := context.Args().Get(0)
		stopContainer(containerId)
		return nil
	},
}

var removeCommand = cli.Command{
	Name:  "rm",
	Usage: "remove unused containers,e.g. mydocker rm 1234567890",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "f", // 强制删除
			Usage: "force delete running container,",
		}},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container id")
		}
		containerId := context.Args().Get(0)
		force := context.Bool("f")
		removeContainer(containerId, force)
		return nil
	},
}
