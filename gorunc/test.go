package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1"
	"log"
	"time"
)

func main() {
	gopts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	addr := "unix:///run/containerd/containerd.sock"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, gopts...)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	req := &v1.VersionRequest{}
	rsp := &v1.VersionResponse{}
	err = conn.Invoke(ctx, "/runtime.v1.RuntimeService/Version", req, rsp)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rsp)

	//v1.PodSandboxConfig{}
	//
	//v1.ContainerConfig{}

}
