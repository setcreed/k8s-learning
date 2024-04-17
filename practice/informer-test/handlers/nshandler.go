package handlers

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

type NsHandler struct {
}

func (ns *NsHandler) OnAdd(obj interface{}, isInInitialList bool) {
	fmt.Printf("namespace added: %s, isInInitialList:%v\n", obj.(*v1.Namespace).Name, isInInitialList)
}

func (ns *NsHandler) OnUpdate(oldObj, newObj interface{}) {
	fmt.Printf("namespace updated: %s\n", newObj.(*v1.Namespace).Name)
}

func (ns *NsHandler) OnDelete(obj interface{}) {
	fmt.Printf("namespace deleted: %s\n", obj.(*v1.Namespace).Name)
}
