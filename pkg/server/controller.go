//go:generate protoc -I ../../internal/pb --go_out=plugins=grpc:../../internal/pb ../../internal/pb/k8s_service.proto

package server

import (
	"context"
	"encoding/json"
	"github.com/Tlantic/k8s-sidecar/internal/manager"
	"github.com/Tlantic/k8s-sidecar/internal/pb"
	"k8s.io/api/batch/v1beta1"
)

type K8sService struct {
}

func NewK8sService() *K8sService {
	return &K8sService{}
}

func (s *K8sService) GetConfigMap(ctx context.Context, in *pb.GetConfigMapRequest) (*pb.GetConfigMapResponse, error) {
	cronManager, err := manager.NewKube(&manager.KubeManagerOptions{
		Config:    in.Kubeconfig,
		Namespace: in.Namespace,
		Timeout:   10,
	})
	if err != nil {
		return nil, err
	}
	data, err := cronManager.GetConfigMap(in.Namespace, in.Key)

	if err != nil {
		return &pb.GetConfigMapResponse{}, err
	}

	return &pb.GetConfigMapResponse{
		Config: data.Data[in.Key],
	}, nil
}

func (s *K8sService) GetCronJobs(ctx context.Context, in *pb.GetCronJobsRequest) (*pb.GetCronJobsResponse, error) {
	cronManager, err := manager.NewKube(&manager.KubeManagerOptions{
		Config:    in.Kubeconfig,
		Namespace: in.Namespace,
		Timeout:   10,
	})
	if err != nil {
		return nil, err
	}
	list, err := cronManager.ListCronJobs()
	if err != nil {
		return nil, err
	}

	cronJobs := make([]*pb.CronJob, len(list.Items))
	for index, item := range list.Items {
		cronJobs[index] = &pb.CronJob{
			Name: item.Name,
		}
	}

	return &pb.GetCronJobsResponse{
		CronJobs: cronJobs,
	}, nil
}

func (s *K8sService) GetCronJob(ctx context.Context, in *pb.GetCronJobRequest) (*pb.GetCronJobResponse, error) {
	cronManager, err := manager.NewKube(&manager.KubeManagerOptions{
		Config:    in.Kubeconfig,
		Namespace: in.Namespace,
		Timeout:   10,
	})
	if err != nil {
		return nil, err
	}
	cronJob, err := cronManager.GetCronJob(in.Id)
	return &pb.GetCronJobResponse{
		CronJob: &pb.CronJob{
			Name: cronJob.Name,
		},
	}, err
}

func (s *K8sService) CreateCronJob(ctx context.Context, in *pb.CreateCronJobsRequest) (*pb.CreateCronJobsResponse, error) {
	cronManager, err := manager.NewKube(&manager.KubeManagerOptions{
		Config:    in.Kubeconfig,
		Namespace: in.Namespace,
		Timeout:   10,
	})
	if err != nil {
		return nil, err
	}
	var jobTemplateData v1beta1.CronJob
	err = json.Unmarshal([]byte(in.Template), &jobTemplateData)
	if err != nil {
		return nil, err
	}

	err = cronManager.CreateCronJob(&jobTemplateData, false)
	return &pb.CreateCronJobsResponse{}, err
}

func (s *K8sService) DeleteCronJob(ctx context.Context, in *pb.DeleteCronJobsRequest) (*pb.DeleteCronJobsResponse, error) {
	cronManager, err := manager.NewKube(&manager.KubeManagerOptions{
		Config:    in.Kubeconfig,
		Namespace: in.Namespace,
		Timeout:   10,
	})
	if err != nil {
		return nil, err
	}
	err = cronManager.DeleteCronJob(in.Name)
	return &pb.DeleteCronJobsResponse{}, err
}
