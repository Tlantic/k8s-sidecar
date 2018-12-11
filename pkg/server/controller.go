//go:generate protoc -I ../../internal/pb --go_out=plugins=grpc:../../internal/pb ../../internal/pb/k8s_service.proto

package server

import (
	"context"
	"encoding/json"
	"github.com/Tlantic/k8s-sidecar/internal/manager"
	"github.com/Tlantic/k8s-sidecar/internal/pb"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
)

type K8sService struct {
	manager *manager.KubeManager
}

func NewK8sService(manager *manager.KubeManager) *K8sService {
	return &K8sService{manager}
}

func (s *K8sService) GetConfigMap(ctx context.Context, in *pb.GetConfigMapRequest) (*pb.GetConfigMapResponse, error) {
	data, err := s.manager.GetConfigMap(in.Key)

	if err != nil {
		return &pb.GetConfigMapResponse{}, err
	}

	return &pb.GetConfigMapResponse{
		Config: data.Data[in.Key],
	}, nil
}

func (s *K8sService) GetCronJobs(ctx context.Context, in *pb.GetCronJobsRequest) (*pb.GetCronJobsResponse, error) {
	list, err := s.manager.ListCronJobs()
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
	cronJob, err := s.manager.GetCronJob(in.Id)
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

	err = s.manager.CreateCronJob(&jobTemplateData, false)
	return &pb.CreateCronJobsResponse{}, err
}

func (s *K8sService) DeleteCronJob(ctx context.Context, in *pb.DeleteCronJobsRequest) (*pb.DeleteCronJobsResponse, error) {
	err := s.manager.DeleteCronJob(in.Name)
	return &pb.DeleteCronJobsResponse{}, err
}

func (s *K8sService) CreateJob(ctx context.Context, in *pb.CreateJobRequest) (*pb.CreateJobResponse, error) {
	var jobTemplateData batchv1.Job
	err := json.Unmarshal([]byte(in.Template), &jobTemplateData)
	if err != nil {
		return nil, err
	}

	err = s.manager.CreateJob(&jobTemplateData, false)
	return &pb.CreateJobResponse{}, err
}

func (s *K8sService) DeleteJob(ctx context.Context, in *pb.DeleteJobRequest) (*pb.DeleteJobResponse, error) {
	err := s.manager.DeleteJob(in.Name)
	return &pb.DeleteJobResponse{}, err
}
