package k8s

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

// Kubeclient is a simplified Kubernetes client
type Kubeclient struct {
	client *client.Client
}

// NewClient creates a new kubernetes client using the config file
func NewClient(config string) (*Kubeclient, error) {

	// c, err := clientFromFile("/Users/pablomercado/.kube/config")

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

	return &Kubeclient{client}, nil
}

// Namespaces return all kubernetes namespaces
func (k *Kubeclient) Namespaces() ([]api.Namespace, error) {

	n, err := k.client.Namespaces().List(api.ListOptions{})
	if err != nil {
		return nil, err
	}

	return n.Items, nil
}

// Pods return pods
func (k *Kubeclient) Pods(namespace string) ([]api.Pod, error) {

	p, err := k.client.Pods(namespace).List(api.ListOptions{})
	if err != nil {
		return nil, err
	}

	return p.Items, nil
}

// ReplicationControllers return replication controllers
func (k *Kubeclient) ReplicationControllers(namespace string) ([]api.ReplicationController, error) {

	rc, err := k.client.ReplicationControllers(namespace).List(api.ListOptions{})
	if err != nil {
		return nil, err
	}

	return rc.Items, nil

}

// Deployments return kubernetes deployments
func (k *Kubeclient) Deployments(namespace string) ([]extensions.Deployment, error) {

	ds, err := k.client.Extensions().Deployments(namespace).List(api.ListOptions{})
	if err != nil {
		return nil, err
	}

	return ds.Items, nil
}
