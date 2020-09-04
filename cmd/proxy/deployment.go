package main

import (
	context "context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// import(
// "github.com/ericchiang/k8s"
// appsv1 "github.com/ericchiang/k8s/apis/apps/v1"
// corev1 "github.com/ericchiang/k8s/apis/core/v1"
// )

func int32P(i int32) *int32 {
	return &i
}

func str(s string) *string {
	return &s
}

func labels(name string) map[string]string {
	return map[string]string{
		"name": name,
		"app":  fmt.Sprintf("cscaler-%s", name),
	}
}

func newDeployment(ctx context.Context, namespace, name, image string) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels(name),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels(name),
			},
			Replicas: int32P(1),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels(name),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image:           image,
							Name:            name,
							ImagePullPolicy: "Always",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "PORT",
									Value: "8080",
								},
							},
						},
					},
				},
			},
		},
	}

	return deployment
}
