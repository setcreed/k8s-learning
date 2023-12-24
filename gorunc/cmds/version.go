package cmds

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(c *cobra.Command, args []string) {
		req := &v1.VersionRequest{}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		rsp, err := NewRuntimeService().Version(ctx, req)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Version:", rsp.Version)
		fmt.Println("RuntimeName:", rsp.RuntimeName)
		fmt.Println("RuntimeVersion:", rsp.RuntimeVersion)
		fmt.Println("RuntimeApiVersion:", rsp.RuntimeApiVersion)
	},
}
