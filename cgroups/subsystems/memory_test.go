package subsystems

import (
	"os"
	"path"
	"testing"
	"time"
)

func TestMemoryCgroup(t *testing.T) {
	memSubSys := MemorySubSystem{}
	resConfig := &ResourceConfig{
		MemoryLimit: "1000m",
	}
	testCgroup := "testmemlimit"

	if err := memSubSys.Set(testCgroup, resConfig); err != nil {
		t.Fatalf("cgroup fail %v", err)
	}
	stat, _ := os.Stat(path.Join(findCgroupMountpoint("memory"), testCgroup))
	t.Logf("cgroup stats: %+v", stat)

	if err := memSubSys.Apply(testCgroup, os.Getpid(), resConfig); err != nil {
		t.Fatalf("cgroup Apply %v", err)
	}
	// 将进程移回到根Cgroup节点
	if err := memSubSys.Apply("", os.Getpid(), resConfig); err != nil {
		t.Fatalf("cgroup Apply %v", err)
	}
	time.Sleep(time.Second * 25)
	if err := memSubSys.Remove(testCgroup); err != nil {
		t.Fatalf("cgroup remove %v", err)
	}
}
