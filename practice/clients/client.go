package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/metadata"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/scale"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/setcreed/practice/helper"
)

const (
	DefaultNamespace = "default"
	DefaultObject    = "ngx"
)

var gvr = schema.GroupVersionResource{
	Group:    "apps",
	Version:  "v1",
	Resource: "deployments",
}

func main() {
	restConfig := helper.BuildClientConfig("")

	// clientset，操作k8s内置资源
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}
	dep, err := clientset.AppsV1().Deployments(DefaultNamespace).Get(context.TODO(), DefaultObject, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found deployment %s in namespace %s\n", dep.Name, DefaultNamespace)

	// MetadataClient 获取k8s元数据
	metadataClient, err := metadata.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}
	meta, err := metadataClient.Resource(gvr).Namespace(DefaultNamespace).Get(context.TODO(), DefaultObject, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Metadata: %v\n", meta)

	// dynamicClient，操作k8s自定义资源
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}
	depResource := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
	unstructured, err := dynamicClient.Resource(depResource).Namespace(DefaultNamespace).Get(context.TODO(), DefaultObject, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Dynamic: %v\n", unstructured)
	var appDep appsv1.Deployment
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured.Object, &appDep)
	fmt.Printf("Dynamic: %v\n", appDep.Name)

	// discoveryClient
	discoveryClient := discovery.NewDiscoveryClientForConfigOrDie(restConfig)
	// 获取支持的 API 组和版本
	apiGroups, _, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		panic(err)
	}
	for _, group := range apiGroups {
		fmt.Printf("Group: %s\n", group.Name)
		for _, version := range group.Versions {
			fmt.Printf("Version: %s\n", version.Version)
		}
	}

	// scaleClient
	groupResources, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		panic(err)
	}
	mapper := restmapper.NewDiscoveryRESTMapper(groupResources)
	scaleResolver := scale.NewDiscoveryScaleKindResolver(discoveryClient)
	scalesGetter, err := scale.NewForConfig(restConfig, mapper, dynamic.LegacyAPIPathResolverFunc, scaleResolver)
	if err != nil {
		panic(err)
	}
	gr := schema.GroupResource{Group: "apps", Resource: "deployments"}
	sca, err := scalesGetter.Scales(DefaultNamespace).Get(context.TODO(), gr, "ngx", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	sca.Spec.Replicas = 5
	//_, err = scalesGetter.Scales(DefaultNamespace).Update(context.TODO(), gr, sca, metav1.UpdateOptions{})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Updated replicas of %s to %d\n", "ngx", sca.Spec.Replicas)

	// metricsClient
	metricsClient, err := metricsv.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}
	podMetricsList, err := metricsClient.MetricsV1beta1().PodMetricses(DefaultNamespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=nginx",
	})
	if err != nil {
		panic(err)
	}
	for _, podMetrics := range podMetricsList.Items {
		// 查看每个 Container 的 CPU 和内存使用情况
		for _, container := range podMetrics.Containers {
			fmt.Printf("Container %s: CPU %v, Memory %v\n", container.Name, container.Usage.Cpu(), container.Usage.Memory())
		}
	}
}
