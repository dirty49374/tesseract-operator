package outgoingportal

import (
	tesseractv1alpha1 "github.com/dirty49374/tesseract-operator/pkg/apis/tesseract/v1alpha1"
	"github.com/dirty49374/tesseract-operator/pkg/certs"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newSecretForCR(cr *tesseractv1alpha1.OutgoingPortal, certs *certs.Certs) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-outgoing-portal",
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app":      cr.Name + "-outgoing-portal",
				"heritage": "tesseract",
			},
		},
		Data: map[string][]byte{
			"ca.crt":     []byte(certs.TrustedCa),
			"client.crt": []byte(certs.Certificate),
			"client.key": []byte(certs.PrivateKey),
		},
	}
}
