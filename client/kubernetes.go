package client

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

// Kubernetes is a simplified client
type Kubernetes struct {
	client *client.Client
}

// NewKubernetesClient creates a new kubernetes client using the config file
func NewKubernetesClient(config string) (*Kubernetes, error) {

	cfgFile, err := clientcmd.LoadFromFile(config)
	if err != nil {
		return nil, fmt.Errorf("error while loading kubeconfig from file %v: %v", config, err)
	}

	cfg, err := clientcmd.NewDefaultClientConfig(*cfgFile, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error while creating kubeconfig: %v", err)
	}

	client, err := client.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("error while creating client: %v", err)
	}

	return &Kubernetes{client}, nil
}

// Namespaces return all kubernetes namespaces
func (k *Kubernetes) Namespaces() ([]api.Namespace, error) {

	n, err := k.client.Namespaces().List(api.ListOptions{})
	if err != nil {
		return nil, err
	}

	return n.Items, nil
}

// Pods return pods
func (k *Kubernetes) Pods(namespace string) ([]api.Pod, error) {

	p, err := k.client.Pods(namespace).List(api.ListOptions{})
	if err != nil {
		return nil, err
	}

	return p.Items, nil
}

// ReplicationControllers return replication controllers
func (k *Kubernetes) ReplicationControllers(namespace string) ([]api.ReplicationController, error) {

	rc, err := k.client.ReplicationControllers(namespace).List(api.ListOptions{})
	if err != nil {
		return nil, err
	}

	return rc.Items, nil

}

// Deployments return kubernetes deployments
func (k *Kubernetes) Deployments(namespace string) ([]extensions.Deployment, error) {

	ds, err := k.client.Extensions().Deployments(namespace).List(api.ListOptions{})
	if err != nil {
		return nil, err
	}

	return ds.Items, nil
}
