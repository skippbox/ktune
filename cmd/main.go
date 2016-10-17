package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	// "github.com/odacremolbap/grisou/image"
	home "github.com/mitchellh/go-homedir"
	"github.com/skippbox/ktune/client"
	"github.com/skippbox/ktune/worker"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var debug bool

//var dockerhubURL string
var kubeconfig string
var frequency int
var namespace string

func init() {

	flag.IntVar(&frequency, "frequency", 10, "Frequency in seconds between deployments/images checks")
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file")
	flag.StringVar(&namespace, "namespace", "default", "Namespace where deployments are located")
	//flag.StringVarP(&dockerhubURL, "dockerhubURL", "u", "", "Dockerhub image URL")
	flag.BoolVarP(&debug, "debug", "d", false, "enable debug messages")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: grisou [OPTIONS] \n")
		fmt.Fprintf(os.Stderr, "\nCanary deployer based on Docker hub image changes\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {

	log.Info("Starting grisou")

	dcc, err := createDeploymentCanaryController()
	if err != nil {
		log.Fatal(err)
	}

	for true {
		log.Infof("Worker iteration at %s", time.Now())
		time.Sleep(time.Duration(frequency) * time.Second)
		dcc.Check()
	}
}

func createDeploymentCanaryController() (*worker.DeploymentCanaryController, error) {

	if kubeconfig == "" {
		h, err := home.Dir()
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get home directory")
		}
		kubeconfig = filepath.Join(h, ".kube", "config")
	}

	k, err := client.NewKubernetesClient(kubeconfig, namespace)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create kubernetes client config")
	}

	d := client.NewDockerHubClient()

	dcc, err := worker.NewDeploymentCanaryController(k, d)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create deployment canary controller")
	}

	return dcc, nil

}
