package main

import (
	"os"
	"os/exec"
)

func main() {
	file, _ := os.Open("./114514.txt")

	cmd := exec.Command("/bin/bash", "-c", "cat <&3")
	cmd.ExtraFiles = []*os.File{file}

	cmd.Stdout = os.Stdout
	cmd.Run()
}
