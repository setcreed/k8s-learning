package cmds

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1"

	"gorunc/utils"
)

// 镜像相关的处理
var imagesCmd = &cobra.Command{
	Use: "images",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		req := &v1.ListImagesRequest{}
		rsp, err := NewImageService().ListImages(ctx, req)
		if err != nil {
			log.Fatalln(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"镜像", "标签", "ID", "大小"})
		for _, img := range rsp.GetImages() {
			imageName, _ := utils.ParseRepoDigest(img.RepoDigests)
			repoTag := utils.ParseRepoTag(img.RepoTags, imageName)[0] // 取到镜像名和标签
			row := []string{imageName, repoTag[1], utils.ParseImageID(img.Id), utils.ParseSize(img.Size_)}
			table.Append(row)
		}
		utils.SetTable(table)
		table.Render()
	},
}
