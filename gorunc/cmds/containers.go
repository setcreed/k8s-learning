package cmds

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	mobyterm "github.com/moby/term"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	restclient "k8s.io/client-go/rest"
	remoteclient "k8s.io/client-go/tools/remotecommand"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1"
	"k8s.io/kubectl/pkg/util/term"

	"github.com/setcreed/gorunc/utils"
)

var containersCmd = &cobra.Command{
	Use:     "run", //单创建 pod
	Example: "run podid container-config.yaml  pod-config.yaml",
	Run: func(c *cobra.Command, args []string) {
		if len(args) < 3 {
			log.Fatalln("参数不完整")
		}
		podId, containConfig, podConfig := "", "", ""
		// 一共三个参数。
		podId = args[0]
		containConfig = args[1]
		podConfig = args[2]

		config := &v1.ContainerConfig{}
		err := utils.YamlFile2Struct(containConfig, config)
		if err != nil {
			log.Fatalln(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		//POD sandbox对应的配置对象
		pConfig := &v1.PodSandboxConfig{}
		err = utils.YamlFile2Struct(podConfig, pConfig)
		if err != nil {
			log.Fatalln(err)
		}
		req := &v1.CreateContainerRequest{
			PodSandboxId:  podId,   //必须要传
			Config:        config,  //容器配置
			SandboxConfig: pConfig, //pod配置 。必须要传
		}

		rsp, err := NewRuntimeService().
			CreateContainer(ctx, req)
		if err != nil {
			log.Fatalln(err)
		}
		_, err = NewRuntimeService().StartContainer(ctx, &v1.StartContainerRequest{
			ContainerId: rsp.ContainerId,
		})
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(rsp.ContainerId) //打印容器ID
	},
}

// 容器列表   类似docker  ps
var containersListCmd = &cobra.Command{
	Use:     "ps", //打印容器
	Example: "ps",
	Run: func(c *cobra.Command, args []string) {

		listReq := &v1.ListContainersRequest{}
		rsp, err := NewRuntimeService().ListContainers(context.Background(), listReq)
		if err != nil {
			log.Fatalln(err)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "名称", "镜像", "状态"})
		for _, c := range rsp.GetContainers() {

			row := []string{utils.ParseContainerID(c.Id), c.Metadata.Name, c.Image.GetImage(),
				strings.Replace(c.State.String(), "CONTAINER_", "", -1)}
			table.Append(row)
		}
		utils.SetTable(table)
		table.Render()
	},
}

// 容器 exec
var containersExecCmd = &cobra.Command{
	Use:     "exec", //打印容器
	Example: "exec",
	Run: func(c *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatalln("error params")
		}
		execReq := &v1.ExecRequest{
			Cmd:         args[1:],
			Stdin:       true,
			Stdout:      true,
			Stderr:      !TTY, // TTY的时候 ，这个值必须是  false
			Tty:         TTY,
			ContainerId: args[0],
		}

		execRsp, err := NewRuntimeService().Exec(context.Background(), execReq)
		if err != nil {
			log.Fatalln(err)
		}
		URL, err := url.Parse(execRsp.Url)
		if err != nil {
			log.Fatalln(err)
		}
		exec, err := remoteclient.NewSPDYExecutor(&restclient.Config{TLSClientConfig: restclient.TLSClientConfig{Insecure: true}}, "POST", URL)

		if !TTY { //非 终端模式
			streamOptions := remoteclient.StreamOptions{
				Stdout: os.Stdout,
				Stderr: os.Stderr,
				Stdin:  os.Stdin,
			}
			err = exec.Stream(streamOptions)
			if err != nil {
				log.Fatalln(err)
			}
			return
		}

		//下面是终端模式
		stdin, stdout, stderr := mobyterm.StdStreams()
		streamOptions := remoteclient.StreamOptions{
			Stdout: stdout,
			Stderr: stderr,
			Stdin:  stdin,
			Tty:    TTY,
		}

		t := term.TTY{
			In:  stdin,
			Out: stdout,
			Raw: true,
		}
		streamOptions.TerminalSizeQueue = t.MonitorSize(t.GetSize())
		err = t.Safe(func() error {
			return exec.Stream(streamOptions)
		})
		if err != nil {
			log.Fatalln(err)
		}
	},
}
