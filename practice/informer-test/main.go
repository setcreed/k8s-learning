package main

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/setcreed/practice/helper"
)

func main() {
	config := helper.BuildClientConfig("")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	// This is the informer
	lw := cache.NewListWatchFromClient(clientset.AppsV1().RESTClient(), "deployments", "default", fields.Everything())
	_, informer := cache.NewInformer(lw, &appsv1.Deployment{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("Deployment added", obj.(*appsv1.Deployment).GetName())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("Deployment update", newObj.(*appsv1.Deployment).GetName())
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("Deployment delete", obj.(*appsv1.Deployment).GetName())
		},
	})

	stop := make(chan struct{})
	defer close(stop)
	go informer.Run(stop)

	select {}
}
