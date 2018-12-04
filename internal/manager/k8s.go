package manager

import (
	"time"

	batchv1 "k8s.io/api/batch/v1"
	batchv1beta "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//KubeCronManager ...
type KubeCronManager struct {
	client    *kubernetes.Clientset
	jobs      map[string]*batchv1.Job
	namespace string
}

type KubeCronManagerOptions struct {
	Config    string
	Namespace string
	Timeout   int
}

//NewKubeCron ...
func NewKubeCron(options *KubeCronManagerOptions) (*KubeCronManager, error) {
	var err error
	k := new(KubeCronManager)
	k.namespace = options.Namespace

	k.client, err = newKubeClientSet(options.Config, options.Timeout)

	if err != nil {
		return nil, err
	}

	return k, nil
}

// NewKubeClientSet creates and initializes a Kubernetes API client to manage our jobs
func newKubeClientSet(kubeconfig string, kubeTimeout int) (*kubernetes.Clientset, error) {
	var err error
	var config *rest.Config

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	} else {

		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
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

/*
* Cronjob Funcs
 */

//GetCronJob ...
func (km *KubeCronManager) GetCronJob(name string) (*batchv1beta.CronJob, error) {
	return km.client.BatchV1beta1().CronJobs(km.namespace).Get(name, metav1.GetOptions{})
}

//CreateCronJob ...
func (km *KubeCronManager) CreateCronJob(job *batchv1beta.CronJob, wait bool) error {
	job.Spec.ConcurrencyPolicy = batchv1beta.ReplaceConcurrent
	if _, err := km.client.BatchV1beta1().CronJobs(km.namespace).Create(job); err != nil {
		return err
	}

	if wait {
		return km.WaitForCronJob(job.Name, km.namespace, 2*time.Minute)
	}

	return nil

}

//DeleteCronJob ...
func (km *KubeCronManager) DeleteCronJob(name string) error {
	policy := metav1.DeletePropagationBackground
	return km.client.BatchV1beta1().CronJobs(km.namespace).Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
}

//ListCronJobs ...
func (km *KubeCronManager) ListCronJobs() (*batchv1beta.CronJobList, error) {
	return km.client.BatchV1beta1().CronJobs(km.namespace).List(metav1.ListOptions{})
}

//WaitForCronJob ...
func (km *KubeCronManager) WaitForCronJob(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		job, err := km.GetCronJob(name)
		if err != nil {
			return false, err
		}
		if len(job.Status.Active) == 0 {
			return false, nil
		}
		return true, nil
	})
}

/*
* Job Funcs
 */

//GetJob ...
func (km *KubeCronManager) GetJob(name string) (*batchv1.Job, error) {
	return km.client.BatchV1().Jobs(km.namespace).Get(name, metav1.GetOptions{})
}

//ListJobs ...
func (km *KubeCronManager) ListJobs() (*batchv1.JobList, error) {
	return km.client.BatchV1().Jobs(km.namespace).List(metav1.ListOptions{})
}

//DeleteJob ...
func (km *KubeCronManager) DeleteJob(name string) error {
	return km.client.BatchV1().Jobs(km.namespace).Delete(name, &metav1.DeleteOptions{})
}

//CreateJob ...
func (km *KubeCronManager) CreateJob(job *batchv1.Job, wait bool) error {
	if _, err := km.client.BatchV1().Jobs(km.namespace).Create(job); err != nil {
		return err
	}

	if wait {
		return km.WaitForJob(job.Name, km.namespace, time.Minute)
	}

	return nil

}

// WaitForJob waits until job deployment has completed
func (km *KubeCronManager) WaitForJob(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		job, err := km.GetJob(name)
		if err != nil {
			return false, err
		}
		if job.Status.Active == 0 && job.Status.Succeeded == 0 {
			return false, nil
		}

		return true, nil
	})
}
