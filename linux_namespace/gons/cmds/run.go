package cmds

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

const (
	SELF   = "/proc/self/exe" // Linux中代表当前的程序
	ALPINE = "/root/alpine"   // 写死的路径，便于测试
)

var runCommand = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd := exec.Command(SELF, "exec", "/bin/sh")
		fmt.Println(runCmd.Args)
		runCmd.SysProcAttr = &syscall.SysProcAttr{
			// 用户隔离
			Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID,
			UidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      os.Getuid(),
					Size:        1,
				},
			},
			GidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      os.Getgid(),
					Size:        1,
				},
			},
		}
		runCmd.Stdin = os.Stdin
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		if err := runCmd.Start(); err != nil {
			log.Fatalln(err)
		}
		runCmd.Wait()
	},
}
