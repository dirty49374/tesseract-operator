package outgoingportal

import (
	"bytes"
	"fmt"
	"text/template"

	tesseractv1alpha1 "github.com/dirty49374/tesseract-operator/pkg/apis/tesseract/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *ReconcileOutgoingPortal) newConfigMapForCR(cr *tesseractv1alpha1.OutgoingPortal) *corev1.ConfigMap {

	if cr == nil {
		return nil
	}

	if cr.Spec.RemotePorts == nil {
		cr.Spec.RemotePorts = []int32{}
	}

	var buf bytes.Buffer

	fmt.Println("XXXXXXX")
	fmt.Println(cr)
	template := template.Must(template.New("envoyConfig").Parse(r.envoyConfig))
	err := template.Execute(&buf, cr)
	if err != nil {
		fmt.Println("==================================================================")
		fmt.Println(err)
		fmt.Println("==================================================================")
		return nil
	}
	fmt.Println("==================================================================")
	fmt.Println(buf.String())
	fmt.Println("==================================================================")

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-portal",
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app":      cr.Name + "-portal",
				"heritage": "tesseract",
			},
		},
		Data: map[string]string{
			"envoy.yaml": buf.String(),
		},
	}
}
