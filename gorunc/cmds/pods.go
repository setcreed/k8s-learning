package cmds

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"

	"github.com/spf13/cobra"
	"gorunc/utils"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1"
)

var podsCmd = &cobra.Command{
	Use: "runp",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalln("请指定POD配置文件")
		}
		config := &v1.PodSandboxConfig{}
		err := utils.YamlFile2Struct(args[0], config)
		if err != nil {
			log.Fatalln(err)
		}
		newUUID := uuid.New()
		config.Metadata.Uid = newUUID.String()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()

		req := &v1.RunPodSandboxRequest{Config: config}
		rsp, err := NewRuntimeService().RunPodSandbox(ctx, req)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(rsp.PodSandboxId)
	},
}
