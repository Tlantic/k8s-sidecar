package server

import (
	"context"
	"github.com/Tlantic/k8s-sidecar/internal/manager"
	"github.com/Tlantic/k8s-sidecar/internal/pb"
	"os"
	"testing"
)

func TestNewK8sService(t *testing.T) {

	kubeManager, err := manager.NewKube(&manager.KubeManagerOptions{
		Config:    os.Getenv("KUBECONFIG"),
		Namespace: os.Getenv("K8S_NAMESPACE"),
		Timeout:   10,
	})
	if err != nil {
		panic(err)
	}

	service := NewK8sService(kubeManager)

	t.Run("GetConfigMap", func(t *testing.T) {
		res, err := service.GetConfigMap(context.Background(), &pb.GetConfigMapRequest{
			Key: "mrs-service-scheduler",
		})
		if err != nil {
			t.Error(err)
		}
		if res == nil {
			t.Errorf("Nil config map")
		}
		t.Log(res)
	})
}
