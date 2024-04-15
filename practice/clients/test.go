package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	syaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"

	"github.com/setcreed/practice/helper"
)

func kustomize(path string) []byte {
	cmd := exec.Command("kustomize", "build", path)
	var ret bytes.Buffer
	cmd.Stdout = &ret
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return nil
	}
	return ret.Bytes()
}

func main() {
	restConfig := helper.BuildClientConfig("")
	dc, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Fatal(err)
	}

	deployYamlStr := kustomize("practice/clients/deploy")
	fmt.Println(string(deployYamlStr))
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(deployYamlStr), len(deployYamlStr))
	for {
		var rawObj runtime.RawExtension
		err := decoder.Decode(&rawObj)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		obj, gvk, err := syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			log.Fatal(err)
		}
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			log.Fatal(err)
		}
		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
		gr, err := restmapper.GetAPIGroupResources(clientset.Discovery())
		if err != nil {
			log.Fatal(err)
		}
		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Fatal(err)
		}

		gvr := schema.GroupVersionResource{
			Group:    mapping.Resource.Group,
			Version:  mapping.Resource.Version,
			Resource: mapping.Resource.Resource,
		}
		_, err = dc.Resource(gvr).Namespace(unstructuredObj.GetNamespace()).Create(context.Background(), unstructuredObj, v1.CreateOptions{})
		if err != nil {
			log.Fatal(err)
		}
		//err = dc.Resource(gvr).Namespace(unstructuredObj.GetNamespace()).Delete(context.Background(), unstructuredObj.GetName(), v1.DeleteOptions{})
		//if err != nil {
		//	log.Fatal(err)
		//}
	}

}
