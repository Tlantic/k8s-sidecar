package manager

import (
	"context"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	informercorev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
	"time"
)

// KubeManager ...
type KubeManager struct {
	client    *kubernetes.Clientset
	jobs      map[string]*batchv1.Job
	namespace string
}

type KubeManagerOptions struct {
	Config    string
	Namespace string
	Timeout   int
}

// NewKube ...
func NewKube(options *KubeManagerOptions) (*KubeManager, error) {
	var err error
	k := new(KubeManager)
	k.namespace = options.Namespace

	k.client, err = newKubeClientSet(options.Config, options.Timeout)

	if err != nil {
		return nil, err
	}

	return k, nil
}

// NewKubeClientSet creates and initializes a Kubernetes API client to manage our jobs
func newKubeClientSet(kubeConfig string, kubeTimeout int) (*kubernetes.Clientset, error) {
	var err error
	var config *rest.Config

	config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}

	if kubeTimeout > 0 {
		config.Timeout = time.Duration(kubeTimeout) * time.Second
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (km *KubeManager) GetConfigMap(ctx context.Context, name string) (*v1.ConfigMap, error) {
	cfgMap, err := km.client.CoreV1().ConfigMaps(km.namespace).Get(ctx, name, metav1.GetOptions{})
	return cfgMap, err

}

func (km *KubeManager) Watch(keys []string, ch chan string, secretInformer informercorev1.ConfigMapInformer) {
	secretInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				data := newObj.(*v1.ConfigMap)
				key, _ := cache.MetaNamespaceKeyFunc(newObj)

				for _, v := range keys {
					if fmt.Sprintf("%s/%s", km.namespace, v) == key {
						ch <- data.Name
					}
				}

			},
		},
	)
}

/*
* Cronjob Funcs
 */

// GetCronJob ...
func (km *KubeManager) GetCronJob(ctx context.Context, name string) (*batchv1.CronJob, error) {
	return km.client.BatchV1().CronJobs(km.namespace).Get(ctx, name, metav1.GetOptions{})
}

// CreateCronJob ...
func (km *KubeManager) CreateCronJob(ctx context.Context, cronJob *batchv1.CronJob, wait bool) error {
	cronJob.Spec.ConcurrencyPolicy = batchv1.ReplaceConcurrent
	if _, err := km.client.BatchV1().CronJobs(km.namespace).Create(ctx, cronJob, metav1.CreateOptions{}); err != nil {
		return err
	}

	if wait {
		return km.WaitForCronJob(ctx, cronJob.Name, km.namespace, 2*time.Minute)
	}

	return nil

}

// DeleteCronJob ...
func (km *KubeManager) DeleteCronJob(ctx context.Context, name string) error {
	policy := metav1.DeletePropagationBackground
	return km.client.BatchV1().CronJobs(km.namespace).Delete(ctx, name, metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
}

// ListCronJobs ...
func (km *KubeManager) ListCronJobs(ctx context.Context) (*batchv1.CronJobList, error) {
	return km.client.BatchV1().CronJobs(km.namespace).List(ctx, metav1.ListOptions{})
}

// WaitForCronJob ...
func (km *KubeManager) WaitForCronJob(ctx context.Context, name, _ string, timeout time.Duration) error {
	return wait.Poll(time.Microsecond*10, timeout, func() (bool, error) {
		job, err := km.GetCronJob(ctx, name)
		if err != nil {
			return false, err
		}
		if strings.Compare(job.Name, name) != 0 {
			return false, nil
		}
		return true, nil
	})
}

/*
* Job Funcs
 */

// GetJob ...
func (km *KubeManager) GetJob(ctx context.Context, name string) (*batchv1.Job, error) {
	return km.client.BatchV1().Jobs(km.namespace).Get(ctx, name, metav1.GetOptions{})
}

// ListJobs ...
func (km *KubeManager) ListJobs(ctx context.Context) (*batchv1.JobList, error) {
	return km.client.BatchV1().Jobs(km.namespace).List(ctx, metav1.ListOptions{})
}

// DeleteJob ...
func (km *KubeManager) DeleteJob(ctx context.Context, name string) error {
	policy := metav1.DeletePropagationBackground
	return km.client.BatchV1().Jobs(km.namespace).Delete(ctx, name, metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
}

// CreateJob ...
func (km *KubeManager) CreateJob(ctx context.Context, job *batchv1.Job, wait bool) error {
	if _, err := km.client.BatchV1().Jobs(km.namespace).Create(ctx, job, metav1.CreateOptions{}); err != nil {
		return err
	}

	if wait {
		return km.WaitForJob(ctx, job.Name, km.namespace, time.Minute)
	}

	return nil

}

// WaitForJob waits until job deployment has completed
func (km *KubeManager) WaitForJob(ctx context.Context, name, _ string, timeout time.Duration) error {
	return wait.Poll(time.Microsecond*5, timeout, func() (bool, error) {
		job, err := km.GetJob(ctx, name)
		if err != nil {
			return false, err
		}
		if job.Status.Active == 0 && job.Status.Succeeded == 0 {
			return false, nil
		}

		return true, nil
	})
}
