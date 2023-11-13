package cmds

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

const ENV = "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

var execCommand = &cobra.Command{
	Use: "exec",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalln("error args")
		}
		runArgs := []string{}
		if len(args) > 1 {
			runArgs = args[1:]
		}

		err := syscall.Chroot(ALPINE) // 设置根目录
		if err != nil {
			log.Fatalln(err)
		}
		err = os.Chdir("/")
		if err != nil {
			log.Fatalln(err)
		}
		syscall.Mount("proc", "/proc", "proc", 0, "") // 挂载

		runCmd := exec.Command(args[0], runArgs...)
		runCmd.Stdin = os.Stdin
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		runCmd.Env = []string{ENV} // 设置环境变量
		if err := runCmd.Start(); err != nil {
			log.Fatalln(err)
		}
		runCmd.Wait()
	},
}
