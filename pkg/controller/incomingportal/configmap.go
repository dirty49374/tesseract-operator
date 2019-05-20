package incomingportal

import (
	"bytes"
	"fmt"

	tesseractv1alpha1 "github.com/dirty49374/tesseract-operator/pkg/apis/tesseract/v1alpha1"
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
			Name:      "tesseract",
			Namespace: tesseractNamespace,
		},
		Data: map[string]string{
			"envoy.yaml": yaml,
		},
	}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *tesseractv1alpha1.IncomingPortal) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
