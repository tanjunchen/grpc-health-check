

package controller

import (
	"log"
	"reflect"
	"time"

	"github.com/tanjunchen/grpc-health/k8s"
	"github.com/tanjunchen/grpc-health/proto"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func generateInfo(svc *corev1.Service, op proto.SyncServiceInfo_Operation) proto.SyncServiceResponse {
	serviceInfo := &proto.SyncServiceInfo{
		Name:              svc.Name,
		ResourceVersion:   svc.ResourceVersion,
		Operation:         op,
		CreationTimeStamp: svc.CreationTimestamp.String(),
		UpdateTimeStamp:   "",
		Labels:            svc.Labels,
		Selector:          svc.Spec.Selector,
	}
	var syncServiceInfo []*proto.SyncServiceInfo
	syncServiceInfo = append(syncServiceInfo, serviceInfo)
	response := proto.SyncServiceResponse{}
	response.SyncServiceInfo = syncServiceInfo
	response.Namespace = svc.Namespace
	return response
}

// Handler interface contains the methods that are required
type ServiceHandler interface {
}

type ServiceController struct {
	K8sClient kubernetes.Interface
	informer  cache.SharedIndexInformer
}

type Condition struct {
	namespace    string
	reSyncPeriod time.Duration
}

func NewServiceController(stopCh <-chan struct{}, informers informers.SharedInformerFactory) (*ServiceController, error) {
	serviceController := ServiceController{}
	serviceController.informer = informers.Core().V1().Services().Informer()
	NewController(stopCh, &serviceController, serviceController.informer)
	return &serviceController, nil
}

func (s *ServiceController) Added(obj interface{}) {
	svc := obj.(*corev1.Service)
	response := generateInfo(svc, 1)
	k8s.MsgChan <- response
	log.Printf("add response: %v", svc)
}

func (s *ServiceController) Updated(newObj interface{}, oldObj interface{}) {
	if oldObj == newObj || reflect.DeepEqual(oldObj, newObj) {
		return
	}
	oldSvc := oldObj.(*corev1.Service)
	newSvc := newObj.(*corev1.Service)
	if oldSvc.ResourceVersion == newSvc.ResourceVersion {
		return
	}
	response := generateInfo(newSvc, 3)
	k8s.MsgChan <- response
	log.Printf("update response: %v", response)
}

func (s *ServiceController) Deleted(obj interface{}) {
	svc := obj.(*corev1.Service)
	response := generateInfo(svc, 2)
	k8s.MsgChan <- response
	log.Printf("delete response: %v", svc)
}
