package k8s

import (
	"log"
	"os"
	"path/filepath"

	"github.com/tanjunchen/grpc-health/proto"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var Client *kubernetes.Clientset

var MsgChan = make(chan proto.SyncServiceResponse, 1024)

// Home returns home path.
func Home() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// Kubeconfig returns kube config path.
func Kubeconfig() string {
	return filepath.Join(Home(), ".kube", "config")
}

func init() {
	clusterConfig, err := clientcmd.BuildConfigFromFlags("", Kubeconfig())
	if err != nil {
		log.Printf("failed to get k8s rest config: %v", err)
	}

	k8s, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		log.Printf("failed to get k8s rest config: %v", err)
	}

	Client = k8s
}
