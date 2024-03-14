package helper

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func BuildClientConfig(filePath string) *rest.Config {
	if len(filePath) == 0 {
		filePath = clientcmd.RecommendedHomeFile
	}
	config, err := clientcmd.BuildConfigFromFlags("", filePath)
	if err != nil {
		panic(err)
	}
	return config
}
