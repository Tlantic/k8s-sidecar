package server

import (
	"context"
	"github.com/Tlantic/k8s-sidecar/internal/pb"
	"os"
	"testing"
)

func TestNewK8sService(t *testing.T) {

	service := NewK8sService()

	t.Run("GetConfigMap", func(t *testing.T) {
		res, err := service.GetConfigMap(context.Background(), &pb.GetConfigMapRequest{
			Kubeconfig: os.Getenv("KUBECONFIG"),
			Namespace:  os.Getenv("K8S_NAMESPACE"),
			Key:        "mrs-service-scheduler",
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
