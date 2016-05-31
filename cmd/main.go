package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	// "github.com/odacremolbap/grisou/image"
	home "github.com/mitchellh/go-homedir"
	"github.com/odacremolbap/grisou/client"
	"github.com/odacremolbap/grisou/cmd/worker"
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

	if kubeconfig == "" {
		h, err := home.Dir()
		if err != nil {
			log.Fatal("Couldn't get home directory. Please, indicate kubeconfig file using flags")
		}
		kubeconfig = filepath.Join(h, ".kube", "config")
	}

	k, err := client.NewKubernetesClient(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	ns, err := k.Namespaces()
	if err != nil {
		log.Fatalf("%v", err)
	}

	found := false
	for _, n := range ns {
		if n.Name == namespace {
			found = true
			break
		}
	}
	if !found {
		log.Fatalf("Namespace '%s' not found", namespace)
	}

	d := client.NewDockerHubClient()

	for true {
		// get kubernetes images
		images, err := worker.UsedImages(k, namespace)
		if err != nil {
			log.Errorf("Failed to get used images from kubernetes: %v", err)
		}
		log.Infof("Images found: %#v", images)

		// get those same images from docker hub
		for _, i := range images {

			t := strings.LastIndex(i, ":")
			// if using latest (no tag) skip version check
			if t == -1 {
				log.Infof("%s is already using latest", i)
				continue
			}

			tag, err := worker.GetLatestTag(d, i[:strings.LastIndex(i, ":")])
			if err != nil {
				log.Errorf("Failed to get latest tag for image '%s': %v", i, err)
			}
			log.Infof("%s latest tag: %#v", i, tag)
		}

		// compare

		// act

		//wait

		time.Sleep(time.Duration(frequency) * time.Second)
	}

}
