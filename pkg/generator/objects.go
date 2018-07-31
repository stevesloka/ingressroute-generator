package generator

import (
	"fmt"

	ingressroutev1 "github.com/stevesloka/ingressroute-generator/pkg/apis/contour/v1beta1"
	clientset "github.com/stevesloka/ingressroute-generator/pkg/generated/clientset/versioned"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func LoopyLoop(num int, namespace string, client *kubernetes.Clientset, contourClient *clientset.Clientset) {
	for i := 0; i < num; i++ {
		svcName := fmt.Sprintf("gensvc-%d", i)
		CreateService(*client, namespace, svcName)
		CreateIngressRoute(*contourClient, namespace, fmt.Sprintf("geningrt%d", i), svcName)
	}
}

func CreateService(client kubernetes.Clientset, namespace, name string) {
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{{
				Name:     "tcp80",
				Port:     80,
				Protocol: v1.ProtocolTCP,
			}},
		},
	}

	_, err := client.CoreV1().Services(namespace).Create(service)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateIngressRoute(client clientset.Clientset, namespace, name, svcname string) {
	route := &ingressroutev1.IngressRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: ingressroutev1.IngressRouteSpec{
			VirtualHost: &ingressroutev1.VirtualHost{
				Fqdn: "r12i.com",
			},
			Routes: []ingressroutev1.Route{{
				Match: "/",
				Services: []ingressroutev1.Service{
					{
						Name: svcname,
						Port: 80,
					},
				},
			}},
		},
	}

	_, err := client.ContourV1beta1().IngressRoutes(namespace).Create(route)
	if err != nil {
		fmt.Println(err)
	}
}
