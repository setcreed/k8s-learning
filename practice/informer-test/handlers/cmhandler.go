package handlers

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

type CmHandler struct {
}

func (c *CmHandler) OnAdd(obj interface{}, isInInitialList bool) {
	fmt.Printf("ConfigMap added: %s, isInInitialList:%v\n", obj.(*v1.ConfigMap).Name, isInInitialList)
}

func (c *CmHandler) OnUpdate(oldObj, newObj interface{}) {

}

func (c *CmHandler) OnDelete(obj interface{}) {
	fmt.Println("ConfigMap deleted:", obj.(*v1.ConfigMap).Name)
}

type CmHandlerNew struct {
}

func (c *CmHandlerNew) OnAdd(obj interface{}, isInInitialList bool) {
	fmt.Printf("CmHandlerNew added: %s, isInInitialList:%v\n", obj.(*v1.ConfigMap).Name, isInInitialList)
}

func (c *CmHandlerNew) OnUpdate(oldObj, newObj interface{}) {

}

func (c *CmHandlerNew) OnDelete(obj interface{}) {
	fmt.Println("CmHandlerNew deleted:", obj.(*v1.ConfigMap).Name)
}
