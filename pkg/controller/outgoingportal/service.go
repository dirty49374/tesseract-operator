package outgoingportal

import (
	"fmt"

	tesseractv1alpha1 "github.com/dirty49374/tesseract-operator/pkg/apis/tesseract/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newServiceForCR(cr *tesseractv1alpha1.OutgoingPortal) *corev1.Service {
	ports := make([]corev1.ServicePort, 0)
	for index, port := range cr.Spec.RemotePorts {
		ports = append(ports, corev1.ServicePort{
			Name:     fmt.Sprintf("port-%d", index),
			Protocol: corev1.ProtocolTCP,
			Port:     port,
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app":      cr.Name + "-outgoing-portal",
				"heritage": "tesseract",
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"app": cr.Name + "-outgoing-portal",
			},
			Ports: ports,
		},
	}
}
