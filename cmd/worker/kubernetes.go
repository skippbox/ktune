package worker

import (
	"fmt"

	"github.com/odacremolbap/grisou/client"
)

// UsedImages return images used in deployments
func UsedImages(client *client.Kubernetes, namespace string) ([]string, error) {

	rcs, err := client.ReplicationControllers(namespace)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get replication controllers: %v", err)
	}

	images := []string{}

	for _, rc := range rcs {

		for _, c := range rc.Spec.Template.Spec.Containers {
			images = append(images, c.Image)
		}
	}

	ds, err := client.Deployments(namespace)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get deployments: %v", err)
	}

	for _, d := range ds {
		for _, c := range d.Spec.Template.Spec.Containers {
			images = append(images, c.Image)
		}
	}

	return images, nil

}
