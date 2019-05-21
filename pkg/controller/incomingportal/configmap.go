package incomingportal

import (
	"bytes"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *ReconcileIncomingPortal) newConfigMapForCR() (*corev1.ConfigMap, error) {

	var buf bytes.Buffer

	err := incomingPortalEnvoyConfig.Execute(&buf, r.portals)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	yaml := buf.String()
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "configmap",
			Namespace: tesseractNamespace,
		},
		Data: map[string]string{
			"envoy.yaml": yaml,
		},
	}, nil
}
