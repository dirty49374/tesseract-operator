package outgoingportal

import (
	"bytes"
	"fmt"

	tesseractv1alpha1 "github.com/dirty49374/tesseract-operator/pkg/apis/tesseract/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *ReconcileOutgoingPortal) newConfigMapForCR(cr *tesseractv1alpha1.OutgoingPortal) (*corev1.ConfigMap, error) {

	if cr == nil {
		return nil, nil
	}

	if cr.Spec.RemotePorts == nil {
		cr.Spec.RemotePorts = []int32{}
	}

	var buf bytes.Buffer

	err := outgoingPortalEnvoyConfig.Execute(&buf, cr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	yaml := buf.String()
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-outgoing-portal",
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app":      cr.Name + "-outgoing-portal",
				"heritage": "tesseract",
			},
		},
		Data: map[string]string{
			"envoy.yaml": yaml,
		},
	}, nil
}
