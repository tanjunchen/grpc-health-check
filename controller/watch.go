package controller

import (
	"log"
	"time"

	"github.com/tanjunchen/grpc-health/k8s"
	v1 "k8s.io/api/core/v1"
	machinery "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
)

func InitSharedInformerFactory() informers.SharedInformerFactory {
	// 构建监听统一过滤标签
	options := func(options *machinery.ListOptions) {
		options.LabelSelector = "app_mesh"
	}
	sharedOptions := []informers.SharedInformerOption{
		informers.WithNamespace(v1.NamespaceAll),
		informers.WithTweakListOptions(options),
	}
	return informers.NewSharedInformerFactoryWithOptions(k8s.Client, time.Second*30, sharedOptions...)
}

func StartWatchController(stopCh <-chan struct{}) {
	informers := InitSharedInformerFactory()

	_, err := NewServiceController(stopCh, informers)
	if err != nil {
		log.Println(err.Error())
	}
}
