package cmds

import (
	"context"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1"
	"log"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const CriAddr = "unix:///run/containerd/containerd.sock"

var grpcClient *grpc.ClientConn

func initClient() {
	gopts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(ctx, CriAddr, gopts...)
	if err != nil {
		log.Fatalln(err)
	}
	grpcClient = conn
}

func NewRuntimeService() v1.RuntimeServiceClient {
	return v1.NewRuntimeServiceClient(grpcClient)
}
func NewImageService() v1.ImageServiceClient {
	return v1.NewImageServiceClient(grpcClient)
}

var TTY bool //终端模式

func RunCmd() {
	cmd := &cobra.Command{
		Use:          "justctl",
		Short:        "jctl",
		Example:      "justctl images",
		SilenceUsage: true,
	}
	initClient()
	containersExecCmd.Flags().BoolVarP(&TTY, "tty", "t", false, "-t")
	cmd.AddCommand(versionCmd, imagesCmd, podsCmd, containersCmd, containersListCmd, containersExecCmd)
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
