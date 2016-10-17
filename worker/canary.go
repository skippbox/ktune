package worker

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/skippbox/ktune/client"
	"github.com/pkg/errors"
)

const ktuneSuffix string = "-ktune"

// DeploymentCanaryController structure to control kubernetes deployments
type DeploymentCanaryController struct {
	kubernetes *client.Kubernetes
	dockerHub  *client.DockerHub
}

// NewDeploymentCanaryController returns a new canary controller for deployments
func NewDeploymentCanaryController(k *client.Kubernetes, d *client.DockerHub) (*DeploymentCanaryController, error) {

	if k == nil {
		return nil, errors.New("Kubernetes client must be set")
	}

	if d == nil {
		return nil, errors.New("Docker Hub client must be set")
	}

	return &DeploymentCanaryController{k, d}, nil
}

// Check creates a canary if a deployment image is outdated
func (dcc *DeploymentCanaryController) Check() error {

	ds, err := dcc.kubernetes.Deployments()
	if err != nil {
		return errors.Wrap(err, "Couldn't retrieve kubernetes deployments")
	}

	dsCanaries := make(map[string]bool, 0)
	for _, d := range ds {
		// if deployment is already a grisou canary, skip
		if strings.HasSuffix(d.Name, ktuneSuffix) {
			ktuned := d.Name[:strings.LastIndex(d.Name, ktuneSuffix)]
			dsCanaries[ktuned] = true
		}
	}

	for _, d := range ds {

		log.Infof("Checking Deployment '%s'", d.Name)

		// if deployment is already a grisou canary, skip
		if strings.HasSuffix(d.Name, ktuneSuffix) {
			log.Debugf("Deployment '%s' is a canary. Skipping", d.Name)
			continue
		}

		// it there is a canary for
		if dsCanaries[d.Name] {
			log.Debugf("Deployment '%s' already has a corresponding canary. Skipping", d.Name)
			continue
		}

		deployCanary := false

		for i := range d.Spec.Template.Spec.Containers {

			// get image tag
			it := strings.Split(d.Spec.Template.Spec.Containers[i].Image, ":")

			if strings.HasPrefix(it[0], "gcr.io") {
				log.Warnf("%s uses gcr.io repository, which is not supported", it[0])
				continue
			}

			if strings.HasPrefix(it[0], "quay.io") {
				log.Warnf("%s uses quay.io repository, which is not supported", it[0])
				continue
			}

			image, err := dcc.dockerHub.GetImageData(it[0])
			if err != nil {
				log.Errorf("couldn't get image data for '%s'", it[0])
			}

			// get image latest tag
			latest := image.GetLatestTag()
			log.Infof("Latest tag found to be '%s' ", latest)
	
			if latest == "" {
				latest = "latest"
			}

			if it[1] == latest {
				log.Debugf("image '%s' is already using the latest version", it)
				continue
			}

			// create canary deployment
			deployCanary = true
			d.Spec.Template.Spec.Containers[i].Image = fmt.Sprintf("%s:%s", it[0], latest)
			log.Infof("image '%s' is outdated. New canary will update to '%s'", it, d.Spec.Template.Spec.Containers[i].Image)
		}

		if deployCanary {
			d.Name = fmt.Sprintf("%s%s", d.Name, ktuneSuffix)
			log.Infof("Creating new deployment '%s' ", d.Name)
			d.Spec.Template.Labels["track"] = "canary"
			d.Spec.Selector = nil
			d.ResourceVersion = ""
			_, err = dcc.kubernetes.CreateDeployment(&d)
			if err != nil {
				log.Errorf("Deployment '%s' could not be created: %v", d.Name, err)
			}
		}
	}

	return nil
}
