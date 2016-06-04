package client

import (
	"github.com/pkg/errors"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

// Kubernetes is a simplified client
type Kubernetes struct {
	client    *client.Client
	namespace string
}

// NewKubernetesClient creates a new kubernetes client using the config file
func NewKubernetesClient(config string, namespace string) (*Kubernetes, error) {

	cfgFile, err := clientcmd.LoadFromFile(config)
	if err != nil {
		return nil, errors.Wrapf(err, "error while loading kubeconfig from file %s", config)
	}

	cfg, err := clientcmd.NewDefaultClientConfig(*cfgFile, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "error creating kubeconfig object")
	}

	client, err := client.New(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating kubernetes client")
	}

	ns, err := client.Namespaces().List(api.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving kubernetes namespaces")
	}

	for _, n := range ns.Items {
		if n.Name == namespace {
			return &Kubernetes{client, namespace}, nil
		}
	}

	return nil, errors.Errorf("Could not find namespace '%s' in kubernetes", namespace)

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
func (k *Kubernetes) Deployments() ([]extensions.Deployment, error) {

	ds, err := k.client.Extensions().Deployments(k.namespace).List(api.ListOptions{})
	if err != nil {
		return nil, err
	}

	return ds.Items, nil
}

// CreateDeployment creates a deployment
func (k *Kubernetes) CreateDeployment(d *extensions.Deployment) (*extensions.Deployment, error) {

	d, err := k.client.Deployments(k.namespace).Create(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}
