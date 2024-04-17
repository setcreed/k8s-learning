package main

import (
	"k8s.io/apimachinery/pkg/fields"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/setcreed/practice/helper"
	"github.com/setcreed/practice/informer-test/handlers"
)

func main() {

	restConfig := helper.BuildClientConfig("")
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Fatal(err)
	}

	// informer基本创建使用
	listWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "configmaps", "default", fields.Everything())
	//_, informer := cache.NewInformer(listWatcher, &v1.ConfigMap{}, 0, &handlers.CmHandler{})
	//informer.Run(wait.NeverStop)

	// sharedInformer创建使用
	sharedInformer := cache.NewSharedInformer(listWatcher, &v1.ConfigMap{}, 0)
	// 同时监听多个
	sharedInformer.AddEventHandler(&handlers.CmHandlerNew{})
	sharedInformer.AddEventHandler(&handlers.CmHandler{})
	sharedInformer.Run(wait.NeverStop)
	select {}
}
