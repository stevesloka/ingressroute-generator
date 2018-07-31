package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/stevesloka/ingressroute-generator/pkg/generator"
	"github.com/stevesloka/ingressroute-generator/pkg/k8s"
)

var (
	kubeClientQPS   float64
	kubeClientBurst int
	numItems        int
)

func init() {
	flag.Float64Var(&kubeClientQPS, "client-qps", 5, "The maximum queries per second (QPS) that can be performed on the Kubernetes API server")
	flag.IntVar(&kubeClientBurst, "client-burst", 10, "The maximum number of queries that can be performed on the Kubernetes API server during a burst")
	flag.IntVar(&numItems, "num-items", 10, "The total number of items to create")
	flag.Parse()
}

func main() {
	log := logrus.New()
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	kubeClient, contourClient := k8s.NewClientWithQPS(*kubeconfig, log, float32(kubeClientQPS), kubeClientBurst)

	generator.LoopyLoop(10, "default", kubeClient, contourClient)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
