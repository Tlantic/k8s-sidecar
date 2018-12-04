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
	CronManager *manager.KubeCronManager
}

func NewK8sService(manager *manager.KubeCronManager) *K8sService {
	return &K8sService{
		CronManager: manager,
	}
}

func (s *K8sService) GetCronJobs(ctx context.Context, in *pb.GetCronJobsRequest) (*pb.GetCronJobsResponse, error) {
	list, err := s.CronManager.ListCronJobs()
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
	cronJob, err := s.CronManager.GetCronJob(in.Id)
	return &pb.GetCronJobResponse{
		CronJob: &pb.CronJob{
			Name: cronJob.Name,
		},
	}, err
}

func (s *K8sService) CreateCronJob(ctx context.Context, in *pb.CreateCronJobsRequest) (*pb.CreateCronJobsResponse, error) {
	var jobTemplateData v1beta1.CronJob
	err := json.Unmarshal([]byte(in.Template), &jobTemplateData)
	if err != nil {
		return nil, err
	}

	err = s.CronManager.CreateCronJob(&jobTemplateData, false)
	return &pb.CreateCronJobsResponse{}, err
}

func (s *K8sService) DeleteCronJob(ctx context.Context, in *pb.DeleteCronJobsRequest) (*pb.DeleteCronJobsResponse, error) {
	err := s.CronManager.DeleteCronJob(in.Name)
	return &pb.DeleteCronJobsResponse{}, err
}
