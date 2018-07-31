package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/stevesloka/ingressroute-generator/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	kubeClientQPS   float64
	kubeClientBurst int
)

func init() {
	flag.Float64Var(&kubeClientQPS, "client-qps", 5, "The maximum queries per second (QPS) that can be performed on the Kubernetes API server")
	flag.IntVar(&kubeClientBurst, "client-burst", 10, "The maximum number of queries that can be performed on the Kubernetes API server during a burst")
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

	kubeClient, err := k8s.NewClientWithQPS(*kubeconfig, log, float32(kubeClientQPS), kubeClientBurst)
	if err != nil {
		log.Fatal("Could not init k8sclient! ", err)
	}

	svcs, _ := kubeClient.Core().Services("default").List(metav1.ListOptions{})
	fmt.Println("svcs: ", len(svcs.Items))

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
