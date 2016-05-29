package main

import (
	"fmt"
	"os"
	"path/filepath"

	// "github.com/odacremolbap/grisou/image"
	home "github.com/mitchellh/go-homedir"
	"github.com/odacremolbap/grisou/k8s"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var debug bool
var dockerhubURL string
var kubeconfig string

func init() {

	flag.BoolVarP(&debug, "debug", "d", false, "enable debug messages")
	flag.StringVarP(&kubeconfig, "kubeconfig", "k", "", "Path to kubeconfig file")
	flag.StringVarP(&dockerhubURL, "dockerhubURL", "u", "", "Dockerhub image URL")

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

	c, err := k8s.NewClient(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	ns, err := c.Namespaces()
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, n := range ns {
		log.Debugf("namespace found '%v'", n.Name)
	}

	rcs, err := c.ReplicationControllers("default")
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, rc := range rcs {
		log.Debugf("rc found '%v'", rc.Name)
	}

	dps, err := c.Deployments("default")
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, dp := range dps {
		log.Debugf("deployment found '%v'", dp.Name)
	}

}
