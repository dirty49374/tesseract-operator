package incomingportal

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tesseract",
			Namespace: tesseractNamespace,
			Labels: map[string]string{
				"app": "tesseract",
			},
			Annotations: map[string]string{
				"external-dns.alpha.kubernetes.io/hostname": tesseractHostname,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeLoadBalancer,
			Selector: map[string]string{
				"app": "tesseract",
			},
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:     "tesseract",
					Protocol: corev1.ProtocolTCP,
					Port:     14443,
				},
			},
		},
	}
}
