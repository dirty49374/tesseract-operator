package incomingportal

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newDeployment(hash string) *appsv1.Deployment {

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tesseract",
			Namespace: tesseractNamespace,
			Labels: map[string]string{
				"app": "tesseract",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &tesseractReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "tesseract",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "tesseract",
					},
					Annotations: map[string]string{
						"config": hash,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "envoy",
							Image: "envoyproxy/envoy:v1.10.0",
							Ports: []corev1.ContainerPort{
								{
									Name:          "tesseract",
									ContainerPort: 14443,
								},
							},
							Command: []string{
								"/usr/local/bin/envoy",
								"-c",
								"/config/envoy.yaml",
								"--log-level",
								"debug",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "tesseract-config",
									MountPath: "/config",
								},
								{
									Name:      "tesseract-secret",
									MountPath: "/secret",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "tesseract-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "tesseract-config",
									},
								},
							},
						},
						{
							Name: "tesseract-secret",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "tesseract-secret",
								},
							},
						},
					},
				},
			},
		},
	}
}
