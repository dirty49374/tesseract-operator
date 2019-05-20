package incomingportal

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	tesseractv1alpha1 "github.com/dirty49374/tesseract-operator/pkg/apis/tesseract/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// const envoyConfig = `
// static_resources:
//   listeners:
//   - address:
//       socket_address:
//         address: 0.0.0.0
//         port_value: 80
//     filter_chains:
//     - filters:
//       - name: envoy.http_connection_manager
//         typed_config:
//           "@type": type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager
//           codec_type: auto
//           stat_prefix: ingress_http
//           access_log:
//           - name: envoy.file_access_log
//             config:
//               path: "/dev/stdout"
//           route_config:
//             name: local_route
//             virtual_hosts:
//             {{- range $name, $ports := . }}
//             {{- range $ports }}
//             - name: "{{ $name }}-{{ . }}"
//               domains:
//               - "{{ $name }}:{{ . }}"
//               {{- if eq . 80 }}
//               - "{{ $name }}"
//               {{- end }}
//               routes:
//               - match:
//                   prefix: /
//                 route:
//                   cluster: "{{ $name }}-{{ . }}"
//             {{- end }}
//             {{- end }}
//           http_filters:
//           - name: envoy.router
//             typed_config: {}
//   clusters:
//   {{- range $name, $ports := . }}
//   {{- range $ports }}
//   - name: "{{ $name }}-{{ . }}"
//     connect_timeout: 0.25s
//     type: strict_dns
//     lb_policy: round_robin
//     load_assignment:
//       cluster_name: "{{ $name }}-{{ . }}"
//       endpoints:
//       - lb_endpoints:
//         - endpoint:
//             address:
//               socket_address:
//                 address: {{ $name }}
//                 port_value: {{ . }}
//   {{- end }}
//   {{- end }}
// admin:
//   access_log_path: "/dev/stdout"
//   address:
//     socket_address:
//       address: 0.0.0.0
//       port_value: 8001
// `

var tesseractNamespace = "sys-tesseract"

func init() {
	if os.Getenv("TESSERACT_NAMESPACE") != "" {
		tesseractNamespace = os.Getenv("TESSERACT_NAMESPACE")
	}
}

func (r *ReconcileIncomingPortal) newConfigMapForCR() *corev1.ConfigMap {

	var buf bytes.Buffer

	template := template.Must(template.New("envoyConfig").Parse(r.envoyConfig))
	err := template.Execute(&buf, r.portals)
	if err != nil {
		fmt.Println("INCOMING ==================================================================")
		fmt.Println(err)
		fmt.Println("INCOMING ==================================================================")
		return nil
	}

	yaml := buf.String()
	fmt.Println("INCOMING ==================================================================")
	fmt.Println(r.portals)
	fmt.Println(yaml)
	fmt.Println("INCOMING ==================================================================")

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tesseract",
			Namespace: tesseractNamespace,
		},
		Data: map[string]string{
			"envoy.yaml": yaml,
		},
	}
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
